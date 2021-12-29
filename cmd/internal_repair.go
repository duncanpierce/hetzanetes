package cmd

import (
	"context"
	"fmt"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/hcloud_client"
	"github.com/duncanpierce/hetzanetes/k8s_client"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

func Repair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			k8sClient := k8s_client.New()
			hcloudClient := hcloud_client.New()
			ticker := time.NewTicker(10 * time.Second)
			for {
				<-ticker.C
				err := repair(k8sClient, hcloudClient, ctx)
				if err != nil {
					log.Printf("error: %s\n", err.Error())
				}
			}
		},
	}
	return cmd
}

func repair(k8sClient *k8s_client.K8sClient, hcloudClient hcloud_client.Client, ctx context.Context) error {
	clusterList, err := k8sClient.GetClusterList()
	if err != nil {
		return err
	}
	if len(clusterList.Items) != 1 {
		return fmt.Errorf("expected 1 Cluster resource but found %d", len(clusterList.Items))
	}
	cluster := clusterList.Items[0]

	servers, err := hcloudClient.Server.AllWithOpts(ctx, hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: label.ClusterNameLabel + "=" + cluster.Name,
		},
	})
	if err != nil {
		return err
	}
	serversInANodeSet := map[*hcloud.Server]bool{}
	for _, nodeSet := range cluster.NodeSets {
		serversInSet := matchServersToNodeSet(servers, cluster.Name+"-"+nodeSet.Name+"-")
		log.Printf("identified %d servers in nodeset %s\n", len(serversInSet), nodeSet.Name)
		for _, s := range serversInSet {
			serversInANodeSet[s] = true
		}
		generationNumber := maxGenerationNumber(serversInSet) + 1
		for i := len(serversInSet); i < nodeSet.Replicas; i++ {
			hcloudClient.CreateServer(env.HCloudToken(), cluster.Name, nodeSet.Name, nodeSet.ApiServer, nodeSet.NodeType, "ubuntu-20.04", generationNumber, cluster.Channel)
			generationNumber++
		}
		for i := nodeSet.Replicas; i < len(serversInSet); i++ {
			lowestGeneration := minGenerationNumber(serversInSet)
			log.Printf("deleting server gen %d from nodeset %s\n", lowestGeneration, nodeSet.Name)
			// TODO taint/drain the node, delete node, then delete server
			hcloudClient.Server.Delete(ctx, serversInSet[lowestGeneration])
		}
	}
	for _, s := range servers {
		if !serversInANodeSet[s] {
			// TODO taint/drain the node, delete node, then delete server
			log.Printf("deleting server %s not in any nodeset\n", s.Name)
			hcloudClient.Server.Delete(ctx, s)
		}
	}
	return err
}

func maxGenerationNumber(serversInSet map[int]*hcloud.Server) int {
	maxGenerationNumber := 0
	for i := range serversInSet {
		if i > maxGenerationNumber {
			maxGenerationNumber = i
		}
	}
	return maxGenerationNumber
}

func minGenerationNumber(serversInSet map[int]*hcloud.Server) int {
	minGenerationNumber := math.MaxInt64
	for i := range serversInSet {
		if i < minGenerationNumber {
			minGenerationNumber = i
		}
	}
	return minGenerationNumber
}

func matchServersToNodeSet(servers []*hcloud.Server, matchingPrefix string) map[int]*hcloud.Server {
	serversInSet := map[int]*hcloud.Server{}
	for _, server := range servers {
		if strings.HasPrefix(server.Name, matchingPrefix) {
			generationText := server.Name[len(matchingPrefix):]
			generationNumber, err := strconv.Atoi(generationText)
			if err == nil {
				serversInSet[generationNumber] = server
			}
			// TODO check if server is attached to the cluster network and switched on
			// TODO check the server is the right server type
			// TODO check the server is/isn't API server
		}
	}
	return serversInSet
}
