package model

import "log"

func (n *NodeStatuses) AddNode(node *NodeStatus) {
	log.Printf("Adding node %s to status", node.Name)
	*n = append(*n, node)
}
