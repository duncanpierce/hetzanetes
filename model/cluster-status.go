package model

import "log"

func (c *ClusterStatus) UpdateVersionRanges() {
	log.Printf("Updating version ranges\n")
	api := VersionRange{}
	worker := VersionRange{}

	for _, nodeSet := range c.NodeSetStatuses {
		api = api.MergeRange(nodeSet.Find(IsApiServer(true), PhaseUpTo(Deleting)).GetVersionRange())
		worker = worker.MergeRange(nodeSet.Find(IsApiServer(false), PhaseUpTo(Deleting)).GetVersionRange())
	}
	c.Versions.Api = api
	c.Versions.Workers = worker
	c.Versions.Nodes = api.MergeRange(worker)
	log.Printf("Version ranges to %#v\n", c.Versions)
}
