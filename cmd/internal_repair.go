package cmd

import (
	"context"
	"github.com/duncanpierce/hetzanetes/hcloud_client"
	"github.com/duncanpierce/hetzanetes/k8s_client"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type (
	Servers []Server
	Server  *hcloud.Server
)

func Repair() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair",
		Short: "Repair the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			k8sClient := k8s_client.New()
			hcloudClient := hcloud_client.New()
			for range time.NewTicker(10 * time.Second).C {
				clusterList, err := k8sClient.GetClusterList()
				if err != nil {
					log.Printf("error getting clusters: %s\n", err.Error())
					continue
				}
				if len(clusterList.Items) != 1 {
					log.Printf("expected 1 Cluster resource but found %d", len(clusterList.Items))
					continue
				}
				cluster := clusterList.Items[0]
				servers, err := hcloudClient.Server.AllWithOpts(ctx, hcloud.ServerListOpts{
					ListOpts: hcloud.ListOpts{
						LabelSelector: label.ClusterNameLabel + "=" + cluster.Name,
					},
				})
				if err != nil {
					log.Printf("error getting servers: %s\n", err.Error())
					continue
				}
				cluster.SetServers(servers)
				err = cluster.Repair(hcloudClient)
				if err != nil {
					log.Printf("error repairing cluster: %s\n", err.Error())
					continue
				}
			}
			return nil
		},
	}
	return cmd
}
