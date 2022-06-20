package model

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/catch"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/hcloud_client"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"log"
	"sort"
	"strconv"
	"strings"
)

type (
	Clusters []*Cluster
	Cluster  struct {
		*Metadata        `json:"metadata,omitempty"`
		*Spec            `json:"spec,omitempty"`
		*Status          `json:"status,omitempty"`
		UnmatchedServers []*hcloud.Server `json:"-"`
	}
	Metadata struct {
		Name string `json:"name,omitempty"`
	}
	Spec struct {
		*Versions `json:"versions,omitempty"`
		NodeSets  `json:"nodeSets,omitempty"`
	}
	Versions struct {
		BaseImage  string `json:"BaseImage,omitempty"`
		Kubernetes string `json:"kubernetes,omitempty"`
		Hetzanetes string `json:"Hetzanetes,omitempty"`
	}
	Status struct {
		*NodeSetStatuses `json:"nodeSets,omitempty"`
	}
	NodeSets []*NodeSet
	NodeSet  struct {
		Name      string   `json:"name"`
		ApiServer bool     `json:"apiServer"`
		Replicas  int      `json:"replicas"`
		NodeType  string   `json:"nodeType"`
		Locations []string `json:"locations,omitempty"`
		Servers   Servers  `json:"-"`
	}
	NodeSetStatuses []*NodeSetStatus
	NodeSetStatus   struct {
		Name       string `json:"name"`
		Generation int    `json:"generation,omitempty"`
	}
	Servers []*Server
	Server  struct {
		*hcloud.Server
		Generation int
		ApiServer  bool
	}
)

func (v *Versions) GetBaseImage() string {
	if v == nil || v.BaseImage == "" {
		return "ubuntu-22.04"
	}
	return v.BaseImage
}

func (v *Versions) GetKubernetes() string {
	if v == nil || v.Kubernetes == "" {
		return "stable"
	}
	return v.Kubernetes
}

func (v *Versions) GetHetzanetes() string {
	if v == nil || v.Hetzanetes == "" {
		return "latest"
	}
	return v.Hetzanetes
}

func (c *Cluster) SetServers(servers []*hcloud.Server) {
matchServers:
	for _, server := range servers {
		for _, nodeSet := range c.NodeSets {
			matchingPrefix := c.Name + "-" + nodeSet.Name + "-"
			if strings.HasPrefix(server.Name, matchingPrefix) {
				generationText := server.Name[len(matchingPrefix):]
				generationNumber, err := strconv.Atoi(generationText)
				if err == nil {
					nodeSet.Servers = append(nodeSet.Servers, &Server{
						Server:     server,
						Generation: generationNumber,
						ApiServer:  nodeSet.ApiServer,
					})
					continue matchServers
				}
			}
		}
		c.UnmatchedServers = append(c.UnmatchedServers, server)
	}
	for _, nodeSet := range c.NodeSets {
		nodeSet.Servers.SortInPlace()
	}
	return
}

func (c *Cluster) Servers() (result Servers) {
	for _, nodeSet := range c.NodeSets {
		result = append(result, nodeSet.Servers...)
	}
	return
}

func (c *Cluster) FirstApiServerNodeSet() *NodeSet {
	for _, nodeSet := range c.NodeSets {
		if nodeSet.ApiServer {
			return nodeSet
		}
	}
	return nil
}

func (n *NodeSet) ServerName(clusterName string, generation int) string {
	return fmt.Sprintf("%s-%s-%d", clusterName, n.Name, generation)
}

func (n *NodeSetStatuses) Named(name string) *NodeSetStatus {
	for _, x := range *n {
		if x.Name == name {
			return x
		}
	}
	return nil
}

func (c *Cluster) NewestApiServer() (newest *Server) {
	servers := c.Servers()
	if len(servers) == 0 {
		return nil
	}
	newest = servers[0]
	for _, server := range servers[1:] {
		if server.ApiServer && server.Created.After(newest.Created) {
			newest = server
		}
	}
	return
}

func (c *Cluster) Repair(hcloudClient hcloud_client.Client) error {
	errs := &catch.Errors{}
	serversMissing := false
	for _, nodeSet := range c.NodeSets {
		repair, err := nodeSet.Repair(c, hcloudClient)
		errs.Add(err)
		if repair {
			serversMissing = true
		}
	}
	// TODO before deleting unmatched servers we should check all required nodes are Ready
	if !serversMissing {
		for _, server := range c.UnmatchedServers {
			errs.Add(hcloudClient.DrainAndDeleteServer(server))
		}
	}
	return errs.OrNil()
}

func (n *NodeSet) Repair(cluster *Cluster, hcloudClient hcloud_client.Client) (bool, error) {
	errs := &catch.Errors{}
	serversMissing := false

	nextGenerationNumber := n.Servers.MaxGeneration() + 1
	for i := len(n.Servers); i < n.Replicas; i++ {
		errs.Add(hcloudClient.CreateServer(env.HCloudToken(), cluster.Name, n.Name, n.ApiServer, n.NodeType, cluster.Versions.GetBaseImage(), nextGenerationNumber, cluster.Versions.GetKubernetes()))
		serversMissing = true
		nextGenerationNumber++
	}
	// TODO before deleting excess servers we should check all required nodes are Ready
	for i := 0; i < len(n.Servers)-n.Replicas; i++ {
		serverToDelete := n.Servers[i]
		log.Printf("deleting server %s\n", serverToDelete)
		errs.Add(hcloudClient.DrainAndDeleteServer(serverToDelete.Server))
	}
	return serversMissing, errs.OrNil()
}

func (s Servers) SortInPlace() {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Generation < s[j].Generation
	})
}

func (s Servers) MaxGeneration() int {
	if len(s) == 0 {
		return 0
	}
	return s[len(s)-1].Generation
}
