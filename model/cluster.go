package model

import "github.com/duncanpierce/hetzanetes/env"

func (c *Cluster) Repair(actions Actions) error {
	if c.Status == nil {
		c.Bootstrap(actions)
	}

	for _, nodeSetSpec := range c.Spec.NodeSets {
		(&c.Status.NodeSetStatuses).CreateIfNecessary(nodeSetSpec)
	}

	c.Status.UpdateVersionRanges()

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

func (c *Cluster) Bootstrap(actions Actions) error {
	c.Status = &ClusterStatus{
		Versions: VersionStatus{},
		ClusterNetwork: ClusterNetworkStatus{
			CloudId: env.HCloudNetworkId(),
			IpRange: env.HCloudNetworkIpRange(),
		},
		NodeSetStatuses: NodeSetStatuses{},
	}

	if err := c.Status.Versions.UpdateReleaseChannels(c.Spec.Versions.GetKubernetes(), actions); err != nil {
		return err
	}

	// TODO find the bootstrap node and add it to NodeSetStatuses in Active state

	return nil
}
