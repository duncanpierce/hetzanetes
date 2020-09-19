package main

import (
	"github.com/duncanpierce/hetzanetes/cmd"
	"github.com/hetznercloud/cli/cli"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	newCLI := cli.NewCLI()
	newCLI.ReadConfig()
	newCLI.ReadEnv()
	client := newCLI.Client()
	ctx := newCLI.Context

	var defaultCmd = &cobra.Command{
		Use: "hetzanetes",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	defaultCmd.AddCommand(cmd.List(client, ctx))
	defaultCmd.AddCommand(cmd.Create(client, ctx))

	if err := defaultCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
