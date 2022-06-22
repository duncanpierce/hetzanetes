package model

func (c *ClusterStatus) UpdateVersionRanges() {
	api := VersionRange{}
	worker := VersionRange{}

	for _, nodeSet := range c.NodeSetStatuses {
		api = api.MergeRange(nodeSet.Find(IsApiServer(true), PhaseUpTo(Deleting)).GetVersionRange())
		worker = worker.MergeRange(nodeSet.Find(IsApiServer(false), PhaseUpTo(Deleting)).GetVersionRange())
	}
	c.Versions.Api = api
	c.Versions.Workers = worker
	c.Versions.Nodes = api.MergeRange(worker)
}
