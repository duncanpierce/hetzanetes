package cmd

import (
	"github.com/spf13/cobra"
)

func Internal() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "internal",
		Short:            "Commands for Hetzanetes internal use, usually run inside the Kubernetes cluster",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		Kustomize(),
		Spike(),
		Repair(),
		Net(),
	)
	return cmd
}
