package model

func (c *ClusterStatus) UpdateVersionRanges() {
	api := VersionRange{}
	worker := VersionRange{}

	for _, nodeSet := range c.NodeSetStatuses {
		api = api.MergeRange(nodeSet.Find(IsApiServer(true), PhaseUpTo(Deleting)).GetVersionRange())
		worker = worker.MergeRange(nodeSet.Find(IsApiServer(false), PhaseUpTo(Deleting)).GetVersionRange())
	}
	c.ApiVersions = api
	c.WorkerVersions = worker
	c.NodeVersions = api.MergeRange(worker)
}
