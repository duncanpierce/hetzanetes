package main

import (
	"github.com/duncanpierce/hetzanetes/cmd"
	"github.com/duncanpierce/hetzanetes/impl"
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
		Use: impl.AppName,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	defaultCmd.AddCommand(cmd.List(client, ctx))
	defaultCmd.AddCommand(cmd.Create(client, ctx))
	defaultCmd.AddCommand(cmd.Delete(client, ctx))
	// TODO it would be nice to add hetzner CLI's 'context' command here, since we share the context, but it's package-private
	// TODO implement "repair" which scans the cluster and recreates resources that are missing, according to the cluster manifest - this would be run as a cronjob in the cluster
	// should probably kick unlabelled servers off the private network, update SSH keys on all servers to latest matching, provision any servers that are missing, update incorrect server and network labels

	if err := defaultCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
