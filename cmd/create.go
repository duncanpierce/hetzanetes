package cmd

import (
	"context"
	"fmt"
	"github.com/duncanpierce/hetzanetes/cloudinit"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"net"
)

// TODO add options to --protected, --backups to enable protection and backups
// TODO maybe protected should be the default
func Create(client *hcloud.Client, ctx context.Context, apiToken string) *cobra.Command {
	var clusterName string
	var ipRange net.IPNet
	var labelsMap map[string]string
	var serverType string
	var osImage string

	cmd := &cobra.Command{
		Use:              "create [FLAGS]",
		Short:            "Create a new cluster",
		Long:             "Create a new Hetzanetes cluster in a new private network.",
		Example:          "  hetzanetes create --name=cluster-1",
		TraverseChildren: true,
		Args:             cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			var labels label.Labels = labelsMap
			labels[label.ClusterNameLabel] = clusterName

			clusterConfig := cloudinit.ClusterConfig{
				ApiToken:           apiToken,
				PrivateNetworkName: clusterName,
				PrivateIpRange:     ipRange.String(),
				K3sInstallEnvVars:  "",
				K3sInstallArgs:     "--disable-cloud-controller",
			}
			cloudInit, err := cloudinit.Template(clusterConfig)
			if err != nil {
				return err
			}

			// TODO check for name collisions on network and API server before starting, and also on server and network labels
			// TODO split this out behind a driver interface to allow --dry-run

			subnets := []hcloud.NetworkSubnet{
				{
					Type:        hcloud.NetworkSubnetTypeCloud,
					IPRange:     &ipRange,
					NetworkZone: hcloud.NetworkZoneEUCentral,
					Gateway:     nil,
				},
			}

			// TODO protect this network - it could be difficult to repair if deleted (e.g. server gets a new interface flannel doesn't know about)
			networkLabels := labels.Copy().Mark(label.PrivateNetworkLabel)
			network, _, err := client.Network.Create(ctx, hcloud.NetworkCreateOpts{
				Name:    clusterName,
				IPRange: &ipRange,
				Subnets: subnets,
				Routes:  nil,
				Labels:  networkLabels,
			})
			if err != nil {
				return err
			}
			fmt.Printf("Created network %s (%s)\n", network.Name, network.IPRange.String())

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
			server, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
				Name:       clusterName + "-api-1",
				ServerType: serverType,
				Image:      image,
				SSHKeys:    sshKeys,
				Location:   nil,
				UserData:   cloudInit,
				Labels:     labels.Copy().Mark(label.ApiServerLabel).Mark(label.WorkerLabel), // TODO --segregate-api to remove this and taint the api server (or have repair do it)
				Networks:   []*hcloud.Network{network},
			})
			if err != nil {
				return err
			}
			fmt.Printf("Created server %s in %s\n", server.Server.Name, server.Server.Datacenter.Name)

			_, _, err = client.Network.Update(ctx, network, hcloud.NetworkUpdateOpts{
				Labels: networkLabels.Set(label.EndpointLabel, server.Server.PublicNet.IPv4.IP.String()),
			})
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&clusterName, "name", "", "Cluster name (required)")
	cmd.MarkFlagRequired("name")
	cmd.Flags().IPNetVar(&ipRange, "cluster-ip-range", net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.IPMask{255, 255, 0, 0}}, "Cluster network IP range")
	// TODO remove cluster-ip-range option? make it an attribute of the network provider?
	// TODO allow create-time-only configuration of pod and service IP ranges? might be easier to leave it on defaults
	cmd.Flags().StringToStringVar(&labelsMap, "label", map[string]string{}, "User-defined labels ('key=value') (can be specified multiple times)")
	cmd.Flags().StringVar(&serverType, "server-type", "cx11", "Server type")
	cmd.Flags().StringVar(&osImage, "os-image", "ubuntu-20.04", "Operating system image")

	return cmd
}
