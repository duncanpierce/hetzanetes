package model

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/env"
	"log"
)

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
		nodeSetStatus.Repair(c, actions)
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

func (c *Cluster) Bootstrap(actions Actions) error {
	log.Printf("Initializing status field\n")
	c.Status = &ClusterStatus{
		Versions: VersionStatus{},
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

	log.Printf("Creating nodeSet statuses\n")
	c.CreateNodeSetStatusesIfNecessary()

	bootstrapApiNodeSetName := c.Spec.NodeSets.FirstApiServerNodeSet().Name
	nodeSetStatus := c.Status.NodeSetStatuses.Named(bootstrapApiNodeSetName)
	nodeSetStatus.Generation = 1
	bootstrapNodeName := fmt.Sprintf("%s-%s-%d", c.Metadata.Name, bootstrapApiNodeSetName, 1)
	log.Printf("Looking for bootstrap API server %s\n", bootstrapNodeName)

	k8sVersion := c.Status.Versions.Target // TODO would be safer to read it from the Node resource in K8s (avoids race condition)
	nodeStatus, err := actions.GetServer(bootstrapNodeName, true, k8sVersion)
	if err != nil {
		return err
	}

	log.Printf("Adding bootstrap node %#v to node status\n", nodeStatus)
	nodeSetStatus.NodeStatuses = append(nodeSetStatus.NodeStatuses, nodeStatus)

	return nil
}
