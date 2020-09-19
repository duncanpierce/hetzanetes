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
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List clusters",
		Long:    `List Hetzanetes clusters, by looking for Hetzner private networks they run in.`,
		Example: `  hetzanetes list`,
		Args:    cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			networks, _, _ := client.Network.List(ctx, hcloud.NetworkListOpts{})
			for _, network := range networks {
				if showAll || hasMarkerLabel(network) {
					fmt.Printf("%s\n", network.Name)
				}
			}
		},
	}
	cmd.Flags().BoolVar(&showAll, "all", false, "show all networks, even if they don't have the expected labels")
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
