package cmd

import (
	"github.com/duncanpierce/hetzanetes/client"
	"github.com/spf13/cobra"
)

func RepairAll(c client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "all",
		Short:            "Scan cluster and repair any problems found",
		Long:             "Scan the cluster and bring it to the expected state, deleting and re-creating resources as necessary. Normally run automatically by the cluster itself.",
		Example:          "  hetzanetes repair all --name=cluster-1",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}
