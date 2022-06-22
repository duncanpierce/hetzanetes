package model

func (c *Cluster) Repair(actions Actions) error {
	c.Status.UpdateVersionRanges()

	for _, nodeSetSpec := range c.Spec.NodeSets {
		(&c.Status.NodeSetStatuses).CreateIfNecessary(nodeSetSpec)
	}

	for _, nodeSetStatus := range c.Status.NodeSetStatuses {
		nodeSetStatus.Repair(c, actions)
	}

	// Action phase change for all nodes
	c.Status.Find(MatchAll()).MakeProgress(c, actions)

	c.Status.UpdateVersionRanges()
	return actions.SaveStatus(c.Metadata.Name, c.Status)
}

func (c *Cluster) FirstApiServerNodeSet() *NodeSetSpec {
	for _, nodeSet := range c.Spec.NodeSets {
		if nodeSet.ApiServer {
			return nodeSet
		}
	}
	return nil
}
