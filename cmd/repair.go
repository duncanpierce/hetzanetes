package cmd

import (
	"github.com/duncanpierce/hetzanetes/client"
	"github.com/spf13/cobra"
)

// TODO this exists to document thinking about how provisioning should work
// Provisioning creates private network and a single API server to which all future servers will join
// Cluster starts with a CronJob manifest that executes `hetzanetes repair all` periodically
// This can be cloned as a Job to run faster on first start
// repair all scans the cluster and Hetzner's API to reconcile discrepancies
// TODO consider how to use and check for hcloud.Actions concerning a resource that are still in progress - to avoid duplicating something already started

func Repair(c client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "repair",
		Short:            "Commands for repairing a cluster",
		Long:             "Commands for repairing a cluster. Normally run automatically by the cluster itself.",
		Example:          "  hetzanetes repair all --name=cluster-1",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		RepairAll(c),
		RepairSsh(c),
		RepairFirewall(c),
		RepairNetwork(c),
		RepairServers(c),
		RepairPackages(c),
	)
	return cmd
}
