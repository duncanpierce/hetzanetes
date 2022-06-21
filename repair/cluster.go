package repair

import "github.com/duncanpierce/hetzanetes/model"

func (c *ClusterStatus) Repair(cluster *model.Cluster, actions Actions) {
	c.UpdateVersionRanges()
	for _, nodeSet := range c.NodeSetStatuses {
		nodeSet.Repair(c, cluster, actions)
	}
	c.UpdateVersionRanges()
}

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
