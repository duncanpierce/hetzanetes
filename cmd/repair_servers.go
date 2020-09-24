package cmd

import (
	"context"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func RepairServers(client *hcloud.Client, ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servers",
		Short: "Creates missing servers and reboots those that need it",
		// priority order:
		// 1. create missing API servers or delete oldest(?) excess ones
		// 2. create missing workers or delete oldest(?) excess ones
		// 3. if all servers are present, reboot single API server which has required reboot for longest, provided it is not already rebooting
		// 4. if all servers are present and no API servers require reboot, reboot single worker which has required reboot for longest, provided it is not already rebooting
		// 5. pick single server that has been rebooting for longest, provided it's over 10 minutes, and delete it, then go to 1.
		Example:          "  hetzanetes repair servers --name=cluster-1",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}
