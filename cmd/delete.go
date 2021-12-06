package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/duncanpierce/hetzanetes/catch"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func Delete(client *hcloud.Client, ctx context.Context) *cobra.Command {
	var clusterName string

	cmd := &cobra.Command{
		Use:              "delete [FLAGS]",
		Short:            "Delete a cluster",
		Long:             "Delete a Hetzanetes cluster and all associated resources including servers and networks.",
		Example:          "  hetzanetes delete --name=cluster-1",
		TraverseChildren: true,
		Args:             cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			network, _, err := client.Network.GetByName(ctx, clusterName)
			if err != nil {
				return err
			}
			if network == nil {
				return errors.New(fmt.Sprintf("cluster network %s not found", clusterName))
			}
			if _, labelled := network.Labels[label.PrivateNetworkLabel]; !labelled {
				return errors.New(fmt.Sprintf("private network %s does not have %s label", clusterName, label.PrivateNetworkLabel))
			}
			labelledServers, err := getServers(client, ctx, clusterName, *network)
			if err != nil {
				return err
			}
			// TODO list the unknown servers on the cluster network rather than returning their number
			discrepancy := len(network.Servers) - len(labelledServers)
			if discrepancy > 0 {
				return errors.New(fmt.Sprintf("%d servers without the correct labels are attached to the cluster network", discrepancy))
			}
			// TODO how do we prevent the cluster from auto-repairing while we're deleting it? Maybe delete should be something the cluster does to itself?

			errs := catch.Errors{}
			for _, server := range labelledServers {
				// TODO retry if deletion fails?
				errs.Catch(client.Server.Delete(ctx, server))
			}
			if len(errs) == 0 {
				errs.Catch(client.Network.Delete(ctx, network))
			}

			apiFirewall, _, apiFirewallErr := client.Firewall.GetByName(ctx, clusterName+"-api")
			errs.Add(apiFirewallErr)
			workerFirewall, _, workerFirewallErr := client.Firewall.GetByName(ctx, clusterName+"-worker")
			errs.Add(workerFirewallErr)
			if apiFirewallErr == nil {
				client.Firewall.Delete(ctx, apiFirewall)
			}
			if workerFirewallErr == nil {
				client.Firewall.Delete(ctx, workerFirewall)
			}
			return errs.OrNil()
		},
	}
	cmd.Flags().StringVar(&clusterName, "name", "", "Cluster name (required)")
	cmd.MarkFlagRequired("name")

	return cmd
}

func getServers(client *hcloud.Client, ctx context.Context, clusterName string, network hcloud.Network) ([]*hcloud.Server, error) {
	servers, err := client.Server.AllWithOpts(ctx, hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: label.ClusterNameLabel + "=" + clusterName,
		},
	})
	if err != nil {
		return nil, err
	}
	for _, server := range servers {
		onClusterNetwork := false
		for _, privateNet := range server.PrivateNet {
			if privateNet.Network.ID == network.ID {
				onClusterNetwork = true
			}
		}
		if !onClusterNetwork {
			return nil, errors.New(fmt.Sprintf("server %s is labelled as being part of the cluster but is not connected to the cluster network", server.Name))
		}
	}
	return servers, nil
}
