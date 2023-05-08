package model

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/json"
	"log"
	"math/rand"
	"time"
)

func (n *NodeSetStatuses) Named(name string) *NodeSetStatus {
	for _, nodeSetStatus := range *n {
		if nodeSetStatus.Name == name {
			return nodeSetStatus
		}
	}
	return nil
}

func (n *NodeSetStatuses) CreateIfNecessary(spec *NodeSetSpec) {
	if n.Named(spec.Name) == nil {
		log.Printf("Creating status for node set '%s'\n", spec.Name)
		*n = append(*n, &NodeSetStatus{
			Name:         spec.Name,
			Generation:   0,
			NodeStatuses: NodeStatuses{},
		})
	}
}

func (n *NodeSetStatuses) BootstrapFrom(cluster *Cluster, servers map[string]*NodeStatus, actions Actions) error {
	bootstrapApiServerName, _ := cluster.BootstrapApiServerName()
	for _, nodeSet := range cluster.Spec.NodeSets {
		replica := 1
		nodeStatuses := NodeStatuses{}

		for ; replica <= nodeSet.Replicas; replica++ {
			serverName := GetServerName(cluster.Metadata.Name, nodeSet.Name, replica)
			server, found := servers[serverName]
			if found {
				server.ApiServer = nodeSet.ApiServer
				server.SetPhase(Creating, "bootstrap server starting up")
				if serverName == bootstrapApiServerName {
					server.SetPhase(Install, "bootstrap api server preinstalled")
					server.SetPhase(Active, "bootstrap api server active")
					node, err := actions.GetNode(bootstrapApiServerName)
					if err != nil {
						return err
					}
					server.Version = node.Status.NodeInfo.KubeletVersion
				}
				nodeStatuses = append(nodeStatuses, server)
			}
		}

		*n = append(*n, &NodeSetStatus{
			Name:         nodeSet.Name,
			Generation:   replica,
			NodeStatuses: nodeStatuses,
		})
	}
	return nil
}

func (n *NodeSetStatus) Repair(cluster *Cluster) {
	log.Printf("Repairing node set '%s'\n", n.Name)
	for _, node := range n.NodeStatuses {
		log.Printf("'%s' status: %s\n", node.Name, json.Format(node))
	}

	target := cluster.Spec.NodeSets.Named(n.Name)

	log.Printf("'%s' node set spec has %d replicas\n", target.Name, target.Replicas)

	// Mark for deletion any stuck nodes:
	n.Find(PhaseUpTo(Joining), LongerThan(10*time.Minute)).SetPhase(Delete, "stuck")

	// TODO Mark for replacement any nodes with wrong baseImage or kubernetes version
	// target version is set in cluster status
	// no worker node can have a higher version that the lowest control plane version
	// no control plane node can be more than 1 minor version ahead of the lowest control plane version

	// Drive towards highest version in cluster
	// When all versions in cluster are the same, drive API servers towards cluster target version

	// TODO replace nodes that don't conform to apiServer bool, baseImage, node type or location
	// need a general "Matches" method to compare to NodeSetStatus / ClusterStatus

	// Create nodes to make up any shortfall against target.Replicas
	// excludes nodes marked for replacement
	// TODO implement maxSurge: limit

	apiServers := cluster.Status.Find(InPhase(Active), IsApiServer(true))
	apiServers.SortByRecency()
	log.Printf("Found %d active API nodes of total %d nodes\n", len(apiServers), len(cluster.Status.Find()))

	if len(apiServers) > 0 {
		joinEndpoint := fmt.Sprintf("https://%s:6443", apiServers[0].ClusterIP)

		for i := len(n.Find(PhaseUpTo(Active))); i < target.Replicas; i++ {
			n.Generation++
			node := &NodeStatus{
				Name:         fmt.Sprintf("%s-%s-%d", cluster.Metadata.Name, target.Name, n.Generation),
				ServerType:   target.ServerType,
				Location:     target.Locations[rand.Intn(len(target.Locations))],
				BaseImage:    target.GetImageOrDefault(),
				ApiServer:    target.ApiServer,
				Version:      cluster.Status.Versions.NewNodeVersion(target.ApiServer),
				JoinEndpoint: joinEndpoint,
			}
			node.SetPhase(Create, "extra server required")
			n.AddNode(node)
		}
	}

	// Mark for deletion the oldest ready nodes beyond the number needed for target.Replicas
	// TODO increasing number of replicas causes new nodes to start creating, then decreasing it shortly after deletes Active nodes, rather than the ones that are being created - which are better candidates for deletion
	// need to take care, because moving a creating node to delete phase means an unpredictable set of actions have already taken place
	// NB there is a race for Joining nodes but not others, provided the phase change isn't happening in a separate controller

	// TODO don't delete an API node which has been assigned as the join endpoint for a not-yet-ready node
	// TODO don't delete the last API server

	// TODO it might be better to delete servers that haven't finished joining yet, rather than wait for them to join then delete earlier servers
	// TODO scale down 1 node at a time

	readyNodes := n.Find(InPhase(Active))
	numberOfUnwantedNodes := len(readyNodes) - target.Replicas
	if numberOfUnwantedNodes > 0 {
		readyNodes[:numberOfUnwantedNodes].SetPhase(Delete, "excess server not required")
	}
}
