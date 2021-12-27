package cmd

import (
	"context"
	"fmt"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/hcloud_client"
	"github.com/duncanpierce/hetzanetes/k8s_client"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"log"
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
			ticker := time.NewTicker(1 * time.Minute)
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
	servers, err := hcloudClient.Server.All(ctx)
	if err != nil {
		return err
	}
	cluster := clusterList.Items[0]
	for _, nodeSet := range cluster.NodeSets {
		serversInSet := matchServersToNodeSet(servers, cluster.Name+"-"+nodeSet.Name+"-")
		maxGenerationNumber := maxGenerationNumber(serversInSet)
		for i := len(serversInSet); i < nodeSet.Replicas; i++ {
			maxGenerationNumber++
			hcloudClient.CreateServer(env.HCloudToken(), cluster.Name, nodeSet.Name, nodeSet.ApiServer, nodeSet.NodeType, "ubuntu-20.04", maxGenerationNumber, cluster.Channel)
		}
		for i := nodeSet.Replicas; i < len(serversInSet); i++ {
			// TODO delete lowest generation server
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

func matchServersToNodeSet(servers []*hcloud.Server, matchingPrefix string) map[int]*hcloud.Server {
	serversInSet := map[int]*hcloud.Server{}
	for _, server := range servers {
		if strings.HasPrefix(server.Name, matchingPrefix) {
			generationNumber, err := strconv.Atoi(server.Name[:len(matchingPrefix)])
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
