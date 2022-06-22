package model

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/rand"
	"time"
)

func (n NodeSetStatuses) Named(name string) *NodeSetStatus {
	for _, nodeSetStatus := range n {
		if nodeSetStatus.Name == name {
			return nodeSetStatus
		}
	}
	return nil
}

func (n *NodeSetStatuses) CreateIfNecessary(spec *NodeSetSpec) {
	if n.Named(spec.Name) == nil {
		*n = append(*n, &NodeSetStatus{
			Name:         spec.Name,
			Generation:   0,
			NodeStatuses: NodeStatuses{},
		})
	}
}

func (n *NodeSetStatus) Repair(cluster *Cluster, actions Actions) {
	target := cluster.Spec.NodeSets.Named(n.Name)

	// Mark for deletion any stuck nodes:
	n.Find(PhaseUpTo(Joining), LongerThan(10*time.Minute)).SetPhase(Delete)

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
	joinEndpoint := fmt.Sprintf("https://%s:6443", apiServers[0].ClusterIP)

	for i := len(n.Find(PhaseUpTo(Active))); i < target.Replicas; i++ {
		n.Generation++
		node := NodeStatus{
			Name:         fmt.Sprintf("%s-%s-%d", cluster.Metadata.Name, target.Name, n.Generation),
			ServerType:   target.ServerType,
			Location:     target.Locations[rand.Intn(len(target.Locations))],
			Created:      time.Now(),
			BaseImage:    cluster.Spec.Versions.BaseImage,
			ApiServer:    target.ApiServer,
			Version:      cluster.Status.VersionStatus.NewNodeVersion(target.ApiServer),
			JoinEndpoint: joinEndpoint,
		}
		node.SetPhase(Create)
		n.AddNode(node)
	}

	// Mark for deletion the oldest ready nodes beyond the number needed for target.Replicas
	// TODO increasing number of replicas causes new nodes to start creating, then decreasing it shortly after deletes Active nodes, rather than the ones that are being created - which are better candidates for deletion
	// need to take care, because moving a creating node to delete phase means an unpredictable set of actions have already taken place
	// NB there is a race for Joining nodes but not others, provided the phase change isn't happening in a separate controller

	// TODO don't delete an API node which has been assigned as the join endpoint for a not-yet-ready node

	readyNodes := n.Find(InPhase(Active))
	numberOfUnwantedNodes := len(readyNodes) - target.Replicas
	if numberOfUnwantedNodes > 0 {
		readyNodes[:numberOfUnwantedNodes].SetPhase(Delete)
	}
}
