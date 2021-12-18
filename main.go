package main

import (
	"context"
	"github.com/duncanpierce/hetzanetes/client"
	"github.com/duncanpierce/hetzanetes/cmd"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	hcloudToken := os.Getenv("HCLOUD_TOKEN")
	if hcloudToken == "" {
		panic("Environment variable HCLOUD_TOKEN must contain a Hetzner Cloud API token")
	}
	c := client.Client{
		Client:  hcloud.NewClient(hcloud.WithToken(hcloudToken)),
		Context: context.Background(),
	}

	var defaultCmd = &cobra.Command{
		Use: label.AppName,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	// TODO work out how to read cluster name from command line without using a --flag
	// TODO need to be able to pass a --context arg
	defaultCmd.AddCommand(
		cmd.List(c.Client, c.Context),
		cmd.Create(c, hcloudToken),
		cmd.Grow(c.Client, c.Context, hcloudToken),
		cmd.Delete(c.Client, c.Context),
		cmd.Repair(c.Client, c.Context),
		cmd.Spike(c.Context),
	)
	// TODO it would be nice to add hetzner CLI's 'context' command here, since we share the context, but it's package-private
	// TODO implement "repair" which scans the cluster and recreates resources that are missing, according to the cluster manifest - this would be run as a cronjob in the cluster
	// should probably kick unlabelled servers off the private network, update SSH keys on all servers to latest matching, provision any servers that are missing, update incorrect server and network labels

	if err := defaultCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
