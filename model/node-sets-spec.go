package model

import (
	"github.com/duncanpierce/hetzanetes/label"
)

func (n NodeSetsSpec) Create(cluster *Cluster, actions Actions, privateNetworkId string, sshPublicKey string) error {
	for _, nodeSet := range n {
		for replica := 1; replica <= nodeSet.Replicas; replica++ {
			labels := label.Labels{}
			if nodeSet.ApiServer {
				labels.Mark(label.ApiServer)
			} else {
				labels.Mark(label.Worker)
			}
			labels.Set(label.Cluster, cluster.Metadata.Name)
			labels.Set(label.NodeSet, nodeSet.Name)
			serverName := GetServerName(cluster.Metadata.Name, nodeSet.Name, replica)
			location := nodeSet.Locations[(replica-1)%len(nodeSet.Locations)]
			_, _, err := actions.CreateServer(serverName, nodeSet.ServerType, nodeSet.GetImageOrDefault(), location, sshPublicKey, privateNetworkId, nil, labels)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (n NodeSetsSpec) BootstrapApiServerNodeSet() *NodeSetSpec {
	for _, nodeSet := range n {
		if nodeSet.ApiServer {
			return nodeSet
		}
	}
	return nil
}

func (n NodeSetsSpec) Named(name string) NodeSetSpec {
	for _, nodeSet := range n {
		if nodeSet.Name == name {
			return *nodeSet
		}
	}
	return NodeSetSpec{
		Name:     name,
		Replicas: 0,
	}
}

func (n *NodeSetSpec) GetImageOrDefault() string {
	if n.Image != "" {
		return n.Image
	}
	return "ubuntu-22.04"
}
