package cmd

import "github.com/spf13/cobra"

func List() *cobra.Command {
	return &cobra.Command{
		Use:          "list",
		Short:        "List clusters",
		Long:         `List Hetzanetes clusters, by looking for Hetzner private networks they run in.`,
		Example:      `  hetzanetes list`,
		SilenceUsage: true,
	}
}
