package cmd

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

const (
	markerLabel string = "hetzanetes"
)

func List(client *hcloud.Client, ctx context.Context) *cobra.Command {
	showAll := false
	var verbose bool

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List clusters",
		Long:    `List Hetzanetes clusters, by looking for Hetzner private networks they run in.`,
		Example: `  hetzanetes list`,
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO use NetworkListOpts with a label rather than filtering client-side
			networks, _, _ := client.Network.List(ctx, hcloud.NetworkListOpts{})
			for _, network := range networks {
				if showAll || hasMarkerLabel(network) {
					fmt.Printf("%s\n", network.Name)
					if verbose {
						fmt.Printf("  ip-range %s\n", network.IPRange)
						for _, subnet := range network.Subnets {
							fmt.Printf("  %s subnet, %s zone, ip-range %v, gateway %v\n", subnet.Type, subnet.NetworkZone, subnet.IPRange, subnet.Gateway)
						}
						for _, server := range network.Servers {
							server, _, err := client.Server.GetByID(ctx, server.ID)
							if err != nil {
								return err
							}
							fmt.Printf("  server %s, %s, %s (%s, %dvcpu, %.2fgb ram), %s in %s\n", server.Name, server.PublicNet.IPv4.IP, server.ServerType.Description, server.ServerType.CPUType, server.ServerType.Cores, server.ServerType.Memory, server.Status, server.Datacenter.Description)
						}
					}
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&showAll, "all", false, "show all networks, even if they don't have the expected labels")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "show more detail, including IP range and connected servers")
	return cmd
}

func hasMarkerLabel(network *hcloud.Network) bool {
	for key := range network.Labels {
		if key == markerLabel {
			return true
		}
	}
	return false
}
