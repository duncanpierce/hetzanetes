package cmd

import (
	"context"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func RepairSsh(client *hcloud.Client, ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "ssh",
		Short:            "Bring SSH keys on all servers up to date with Hetzner API",
		Long:             "Bring SSH keys on all servers up to date with the labelled keys in Hetzner's API. Normally run automatically by the cluster itself.",
		Example:          "  hetzanetes repair ssh --name=cluster-1",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}
