package cmd

import (
	"context"
	"github.com/duncanpierce/hetzanetes/impl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func RepairPackages(client *hcloud.Client, ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "packages",
		Short:            "Bring system software packages up to date", // TODO could be done by unattended-upgrades (or might fight with it), could also update Kubernetes distro
		Example:          impl.AppName + "  hetzanetes repair packages --name=cluster-1",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}
