package cmd

import (
	"github.com/duncanpierce/hetzanetes/client"
	"github.com/spf13/cobra"
)

func RepairPackages(c client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "packages",
		Short:            "Bring system software packages up to date", // TODO could be done by unattended-upgrades (or might fight with it), could also update Kubernetes distro
		Example:          "  hetzanetes repair packages --name=cluster-1",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}
