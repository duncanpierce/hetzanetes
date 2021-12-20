package cmd

import (
	"github.com/spf13/cobra"
)

func Internal() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "internal",
		Short:            "Commands for Hetzanetes internal use",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		Kustomize(),
		Grow(),
		Spike(),
		Net(),
		Watch(),
	)
	return cmd
}
