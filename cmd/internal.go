package cmd

import (
	"github.com/duncanpierce/hetzanetes/client"
	"github.com/spf13/cobra"
)

func Internal(c client.Client, hcloudToken string) *cobra.Command {
	cmd := &cobra.Command{
		Use:              "internal",
		Short:            "Commands for Hetzanetes internal use",
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(
		Grow(c, hcloudToken),
		Spike(c),
		Net(c),
		Watch(c),
	)
	return cmd
}
