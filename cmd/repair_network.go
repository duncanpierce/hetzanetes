package cmd

import (
	"context"
	"github.com/duncanpierce/hetzanetes/impl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func RepairNetwork(client *hcloud.Client, ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Bring private network up to date with cluster configuration",
		// TODO need to decide if provisioning should create the network or if it should be a repair task
		// should create and label network, then scan for servers which need to be attached (missing members) or detached (outsiders) - latter may need to be a separate step
		Example:          impl.AppName + "  hetzanetes repair network --name=cluster-1",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}
