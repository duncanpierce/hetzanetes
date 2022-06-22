package model

func (n *NodeStatuses) AddNode(node NodeStatus) {
	*n = append(*n, node)
}
