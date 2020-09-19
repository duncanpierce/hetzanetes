package cmd

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"net"
)

// TODO add options to --protected, --backups to enable protection and backups
// TODO add option to select SSH key, defaulting to first one found
func Create(client *hcloud.Client, ctx context.Context) *cobra.Command {
	var clusterName string
	var ipRange net.IPNet
	var labels map[string]string
	var serverType string
	var osImage string

	cmd := &cobra.Command{
		Use:              "create [FLAGS]",
		Short:            "Create a new cluster",
		Long:             `Create a new Hetzanetes cluster in a new private network.`,
		Example:          `  hetzanetes create --name=cluster-1`,
		TraverseChildren: true,
		Args:             cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			labels[roleLabel] = "cluster"
			labels[clusterLabel] = clusterName

			// TODO do we need to create a subnet, e.g. from 10.0.0.0/16 create subnet 10.0.0.0/24, rather than taking the whole ip range?
			subnets := []hcloud.NetworkSubnet{
				{
					Type:        hcloud.NetworkSubnetTypeCloud,
					IPRange:     &ipRange,
					NetworkZone: hcloud.NetworkZoneEUCentral,
					Gateway:     nil,
				},
			}

			network, _, err := client.Network.Create(ctx, hcloud.NetworkCreateOpts{
				Name:    clusterName,
				IPRange: &ipRange,
				Subnets: subnets,
				Routes:  nil,
				Labels:  labels,
			})
			if err != nil {
				return err
			}
			fmt.Printf("Created network %s\n", network.Name)

			serverType, _, err := client.ServerType.GetByName(ctx, serverType)
			if err != nil {
				return err
			}
			image, _, err := client.Image.GetByName(ctx, osImage)
			if err != nil {
				return err
			}

			// TODO need to select and pass in an SSH key
			// TODO need an option to select datacenter, although it defaults to fsn1-dc14 anyway
			server, _, err := client.Server.Create(ctx, hcloud.ServerCreateOpts{
				Name:       clusterName + "-api-1",
				ServerType: serverType,
				Image:      image,
				SSHKeys:    nil,
				Location:   nil,
				Datacenter: nil,
				UserData:   "",
				Labels:     map[string]string{roleLabel: "api-server", clusterLabel: clusterName},
				Networks:   []*hcloud.Network{network},
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
	cmd.Flags().IPNetVar(&ipRange, "ip-range", net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.IPMask{255, 255, 0, 0}}, "Network IP range")
	cmd.Flags().StringToStringVar(&labels, "label", map[string]string{}, "User-defined labels ('key=value') (can be specified multiple times)")
	cmd.Flags().StringVar(&serverType, "server-type", "cx11", "Server type")
	cmd.Flags().StringVar(&osImage, "os-image", "ubuntu-20.04", "Operating system image")

	return cmd
}
