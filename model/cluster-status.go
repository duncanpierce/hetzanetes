package model

import (
	"github.com/duncanpierce/hetzanetes/json"
	"log"
)

func (c *ClusterStatus) UpdateVersionRanges() {
	log.Printf("Updating version ranges\n")
	api := &VersionRange{}
	worker := &VersionRange{}

	for _, nodeSet := range c.NodeSetStatuses {
		api = api.MergeRange(nodeSet.Find(IsApiServer(true), PhaseUpTo(Deleting)).GetVersionRange())
		worker = worker.MergeRange(nodeSet.Find(IsApiServer(false), PhaseUpTo(Deleting)).GetVersionRange())
	}
	c.Versions.Api = api
	if worker != nil {
		c.Versions.Workers = worker
	} else {
		c.Versions.Workers = api
	}
	c.Versions.Nodes = api.MergeRange(worker)
	log.Printf("Version ranges to %s\n", json.Format(c.Versions))
}
