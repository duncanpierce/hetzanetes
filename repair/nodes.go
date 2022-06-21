package repair

func (n *NodeStatuses) AddNode(node NodeStatus) {
	*n = append(*n, node)
}

func (n NodeStatusRefs) SetPhase(phase Phase) {
	for i := 0; i < len(n); i++ {
		n[i].SetPhase(phase)
	}
}
