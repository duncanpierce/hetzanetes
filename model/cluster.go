package model

import (
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
		Metadata         `json:"metadata"`
		Spec             `json:"spec"`
		UnmatchedServers []*hcloud.Server
	}
	Metadata struct {
		Name string `json:"name"`
	}
	Spec struct {
		Channel  string `json:"channel"`
		NodeSets `json:"nodeSets"`
	}
	NodeSets []*NodeSet
	NodeSet  struct {
		Name      string   `json:"name"`
		ApiServer bool     `json:"apiServer"`
		Replicas  int      `json:"replicas"`
		NodeType  string   `json:"nodeType"`
		Locations []string `json:"locations"`
		Servers   Servers
	}
	Servers []*Server
	Server  struct {
		*hcloud.Server
		Generation int
		ApiServer  bool
	}
)

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

func (c *Cluster) Repair(hcloudClient hcloud_client.Client) (errs catch.Errors) {
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
	return
}

func (n *NodeSet) Repair(cluster *Cluster, hcloudClient hcloud_client.Client) (serversMissing bool, errs catch.Errors) {
	nextGenerationNumber := n.Servers.MaxGeneration() + 1
	for i := len(n.Servers); i < n.Replicas; i++ {
		errs.Add(hcloudClient.CreateServer(env.HCloudToken(), cluster.Name, n.Name, n.ApiServer, n.NodeType, "ubuntu-20.04", nextGenerationNumber, cluster.Channel))
		serversMissing = true
		nextGenerationNumber++
	}
	// TODO before deleting excess servers we should check all required nodes are Ready
	for i := 0; i < len(n.Servers)-n.Replicas; i++ {
		serverToDelete := n.Servers[i]
		log.Printf("deleting server %s\n", serverToDelete)
		errs.Add(hcloudClient.DrainAndDeleteServer(serverToDelete.Server))
	}
	return
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