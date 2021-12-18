package cmd

import (
	"github.com/duncanpierce/hetzanetes/client"
	"github.com/spf13/cobra"
)

func RepairFirewall(c client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "firewall",
		Short:            "Bring firewall on all servers up to date with cluster configuration",
		Example:          "  hetzanetes repair firewall --name=cluster-1",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}
