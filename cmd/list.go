package cmd

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/client"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func List(c client.Client) *cobra.Command {
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
			networks, _, _ := c.Network.List(c, hcloud.NetworkListOpts{})
			for _, network := range networks {
				_, labelled := network.Labels[label.PrivateNetworkLabel]
				if labelled || showAll {
					fmt.Printf("%s (cluster=%t)\n", network.Name, labelled)
					if verbose {
						fmt.Printf("  ip-range %s\n", network.IPRange)
						for _, subnet := range network.Subnets {
							fmt.Printf("  %s subnet, %s zone, ip-range %v, gateway %v\n", subnet.Type, subnet.NetworkZone, subnet.IPRange, subnet.Gateway)
						}
						for _, server := range network.Servers {
							server, _, err := c.Server.GetByID(c, server.ID)
							if err != nil {
								return err
							}
							fmt.Printf("  server %s (%s) %s (%s %dvcpu %.2fgb ram) %s in %s\n", server.Name, server.PublicNet.IPv4.IP, server.ServerType.Description, server.ServerType.CPUType, server.ServerType.Cores, server.ServerType.Memory, server.Status, server.Datacenter.Description)
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
