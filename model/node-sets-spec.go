package model

import "fmt"

func (n NodeSetsSpec) FirstApiServerNodeSet() *NodeSetSpec {
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

func (n *NodeSetSpec) ServerName(clusterName string, generation int) string {
	return fmt.Sprintf("%s-%s-%d", clusterName, n.Name, generation)
}
