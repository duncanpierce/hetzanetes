package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/tmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"os"
)

// TODO this is a temporary command to add a single worker to the cluster
func Grow(client *hcloud.Client, ctx context.Context, apiToken string) *cobra.Command {
	var clusterName string
	var labelsMap map[string]string
	var serverType string
	var osImage string
	var nodeSuffix string

	cmd := &cobra.Command{
		Use:              "grow [FLAGS]",
		Short:            "Add a single worker to an existing cluster",
		Long:             "Add a single worker to an existing cluster",
		Example:          "  hetzanetes grow --name=cluster-1",
		TraverseChildren: true,
		Args:             cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var labels label.Labels = labelsMap
			labels[label.ClusterNameLabel] = clusterName

			network, _, err := client.Network.GetByName(ctx, clusterName)
			if err != nil {
				return err
			}
			_, labelled := network.Labels[label.PrivateNetworkLabel]
			if !labelled || network.Labels[label.ClusterNameLabel] != clusterName {
				return errors.New(fmt.Sprintf("network %s is not labelled as a cluster", network.Name))
			}

			// TODO find next server number instead of hardcoded 1
			nodeName := clusterName + "-worker-" + nodeSuffix
			ipRange := network.Subnets[0].IPRange

			clusterConfig := tmpl.ClusterConfig{
				JoinToken:          os.Getenv("K3S_TOKEN"),
				ApiEndpoint:        os.Getenv("K3S_URL"),
				HetznerApiToken:    apiToken, // from HCLOUD_TOKEN
				PrivateNetworkName: clusterName,
				PrivateIpRange:     ipRange.String(),
			}
			cloudInit := tmpl.Template(clusterConfig, "add-worker.yaml")

			// TODO check for name collisions new server before starting

			serverType, _, err := client.ServerType.GetByName(ctx, serverType)
			if err != nil {
				return err
			}
			image, _, err := client.Image.GetByName(ctx, osImage)
			if err != nil {
				return err
			}

			// TODO allow a label selector to select keys to use (repair will keep it up to date)
			sshKeys, err := client.SSHKey.All(ctx)
			if err != nil {
				return err
			}

			// Hetzner recommend specifying locations rather than datacenters: https://docs.hetzner.cloud/#servers-create-a-server
			// TODO add --regions option
			t := true
			server, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
				Name:             nodeName,
				ServerType:       serverType,
				Image:            image,
				SSHKeys:          sshKeys,
				Location:         nil,
				StartAfterCreate: &t,
				UserData:         cloudInit,
				Labels:           labels.Copy().Mark(label.WorkerLabel),
				Networks:         []*hcloud.Network{network},
			})
			if err != nil {
				return err
			}

			fmt.Printf("Created server %s in %s\n", server.Server.Name, server.Server.Datacenter.Name)
			return nil
		},
	}
	cmd.Flags().StringVar(&clusterName, "name", "", "Cluster name (required)")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&nodeSuffix, "node-suffix", "1", "Final component of new node name - must be unique within cluster")
	cmd.Flags().StringToStringVar(&labelsMap, "label", map[string]string{}, "User-defined labels ('key=value') (can be specified multiple times)")
	cmd.Flags().StringVar(&serverType, "server-type", "cx11", "Server type")
	cmd.Flags().StringVar(&osImage, "os-image", "ubuntu-20.04", "Operating system image")

	return cmd
}
