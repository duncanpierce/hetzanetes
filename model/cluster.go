package model

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/env"
	"log"
)

func (c *Cluster) Create(actions Actions, privateNetworkId string, sshPublicKey string) error {
	return c.Spec.NodeSets.Create(c, actions, privateNetworkId, sshPublicKey)
}

func (c *Cluster) Repair(actions Actions) error {
	if c.Status == nil {
		log.Printf("No status found - bootstrapping\n")
		err := c.Bootstrap(actions)
		if err != nil {
			return err
		}
	}

	c.CreateNodeSetStatusesIfNecessary()
	c.Status.UpdateVersionRanges()

	for _, nodeSetStatus := range c.Status.NodeSetStatuses {
		nodeSetStatus.Repair(c)
	}

	c.Status.Find().MakeProgress(c, actions)

	c.Status.UpdateVersionRanges()
	return actions.SaveStatus(c.Metadata.Name, c.Status)
}

func (c *Cluster) CreateNodeSetStatusesIfNecessary() {
	for _, nodeSetSpec := range c.Spec.NodeSets {
		(&c.Status.NodeSetStatuses).CreateIfNecessary(nodeSetSpec)
	}
}

func (c *Cluster) BootstrapApiServerName() (string, error) {
	bootstrapApiServerNodeSet := c.Spec.NodeSets.BootstrapApiServerNodeSet()
	if bootstrapApiServerNodeSet == nil {
		return "", fmt.Errorf("cluster does not have an API server")
	}
	return GetServerName(c.Metadata.Name, bootstrapApiServerNodeSet.Name, 1), nil
}

func (c *Cluster) Bootstrap(actions Actions) error {
	log.Printf("Initializing cluster %s status field\n", c.Metadata.Name)
	c.Status = &ClusterStatus{
		Versions: &VersionStatus{},
		ClusterNetwork: ClusterNetworkStatus{
			CloudId: env.HCloudNetworkId(),
			IpRange: env.HCloudNetworkIpRange(),
		},
		NodeSetStatuses: NodeSetStatuses{},
	}

	log.Printf("Downloading release channels and setting target Kubernetes version for cluster\n")
	if err := c.Status.Versions.UpdateReleaseChannels(c.Spec.Versions.GetKubernetes(), actions); err != nil {
		return err
	}

	servers, err := actions.GetServers(c.Metadata.Name)
	if err != nil {
		return err
	}

	c.Status.NodeSetStatuses.BootstrapFrom(c, servers)
	return nil
}
