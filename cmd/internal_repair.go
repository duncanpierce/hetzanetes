package cmd

import (
	"github.com/duncanpierce/hetzanetes/model"
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
			act := model.NewClusterClient()
			for range time.Tick(10 * time.Second) {
				clusterList, err := act.GetClusterList()
				if err != nil {
					log.Printf("error getting clusters: %s\n", err.Error())
					continue
				}
				if len(clusterList.Items) != 1 {
					log.Printf("expected 1 Cluster resource but found %d", len(clusterList.Items))
					continue
				}
				cluster := clusterList.Items[0]
				if err != nil {
					log.Printf("error getting servers: %s\n", err.Error())
					continue
				}

				err = cluster.Repair(act)

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
