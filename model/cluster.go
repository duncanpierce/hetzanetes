package model

import (
	"github.com/duncanpierce/hetzanetes/env"
)

func (c *Cluster) Repair(actions Actions) error {
	if c.Status == nil {
		c.Bootstrap(actions)
	}

	c.CreateNodeSetStatusesIfNecessary()
	c.Status.UpdateVersionRanges()

	for _, nodeSetStatus := range c.Status.NodeSetStatuses {
		nodeSetStatus.Repair(c, actions)
	}

	// Action phase change for all nodes
	c.Status.Find(MatchAll()).MakeProgress(c, actions)

	c.Status.UpdateVersionRanges()
	return actions.SaveStatus(c.Metadata.Name, c.Status)
}

func (c *Cluster) CreateNodeSetStatusesIfNecessary() {
	for _, nodeSetSpec := range c.Spec.NodeSets {
		(&c.Status.NodeSetStatuses).CreateIfNecessary(nodeSetSpec)
	}
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
	c.CreateNodeSetStatusesIfNecessary()

	if err := c.Status.Versions.UpdateReleaseChannels(c.Spec.Versions.GetKubernetes(), actions); err != nil {
		return err
	}

	// TODO find the bootstrap node and add it to NodeSetStatuses in Active state - can't be passed in via cloudinit
	//bootstrapNodeName := fmt.Sprintf("%s-%s-%d", c.Metadata.Name, c.Spec.NodeSets.FirstApiServerNodeSet().Name, 1)
	//actions.GetServerId(bootstrapNodeName)
	// might be better to wait until we can install k3s using SSH ?

	return nil
}
