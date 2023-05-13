package cmd

import (
	"github.com/duncanpierce/hetzanetes/kubeconfig"
	"github.com/duncanpierce/hetzanetes/model"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

type (
	Servers []Server
	Server  *hcloud.Server
)

func Repair() *cobra.Command {
	var kubeconfigFilename string

	cmd := &cobra.Command{
		Use:              "repair",
		Short:            "Repair the cluster",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				server      string
				certificate []byte
				token       string
				err         error
			)
			if kubeconfigFilename != "" {
				kubeconfigContents, err := os.ReadFile(kubeconfigFilename)
				if err != nil {
					return err
				}
				server, certificate, token, err = kubeconfig.FromConfig(kubeconfigContents)
				if err != nil {
					return err
				}
			} else {
				server, certificate, token, err = kubeconfig.InCluster()
				if err != nil {
					return err
				}
			}

			actions, err := model.NewClusterActions(server, certificate, token)
			if err != nil {
				return err
			}
			for {
				clusterList, err := actions.GetClusterList()
				if err != nil {
					log.Printf("error getting clusters: %s\n", err.Error())
				} else if len(clusterList.Items) != 1 {
					log.Printf("expected 1 Cluster resource but found %d", len(clusterList.Items))
				} else {
					cluster := clusterList.Items[0]
					err = cluster.Repair(actions)
					if err != nil {
						log.Printf("error repairing cluster: %s\n", err.Error())
					}
				}
				<-time.Tick(10 * time.Second)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&kubeconfigFilename, "kubeconfig", "k", "", "Name of kubeconfig YAML file to use to connect (not required when running in the cluster)")
	return cmd
}
