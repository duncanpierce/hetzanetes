package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/duncanpierce/hetzanetes/impl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
)

func Delete(client *hcloud.Client, ctx context.Context) *cobra.Command {
	var clusterName string

	cmd := &cobra.Command{
		Use:              "delete [FLAGS]",
		Short:            "Delete a cluster",
		Long:             "Delete a Hetzanetes cluster and all associated resources including servers and networks.",
		Example:          impl.AppName + "  hetzanetes delete --name=cluster-1",
		TraverseChildren: true,
		Args:             cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			network, _, err := client.Network.GetByName(ctx, clusterName)
			if err != nil {
				return err
			}
			if network.Labels[impl.RoleLabel] != impl.ClusterRole {
				return errors.New("the network " + clusterName + " does not have the " + impl.RoleLabel + "=" + impl.ClusterRole + " label")
			}
			apiServers, err := getServers(client, ctx, impl.ApiServerRole, clusterName, *network)
			if err != nil {
				return err
			}
			workers, err := getServers(client, ctx, impl.WorkerRole, clusterName, *network)
			if err != nil {
				return err
			}
			// TODO list the unknown servers on the cluster network rather than returning their number
			servers := append(apiServers, workers...)
			discrepancy := len(network.Servers) - len(servers)
			if discrepancy > 0 {
				return errors.New(fmt.Sprintf("%d servers without the correct labels are attached to the cluster network", discrepancy))
			}
			// TODO how do we prevent the cluster from auto-repairing while we're deleting it? Maybe delete should be something the cluster does to itself?

			for _, server := range servers {
				// TODO retry if deletion fails? or keep going and warn at the end?
				client.Server.Delete(ctx, server)
			}
			client.Network.Delete(ctx, network)
			return nil
		},
	}
	cmd.Flags().StringVar(&clusterName, "name", "", "Cluster name (required)")
	cmd.MarkFlagRequired("name")

	return cmd
}

func getServers(client *hcloud.Client, ctx context.Context, role, clusterName string, network hcloud.Network) ([]*hcloud.Server, error) {
	servers, err := client.Server.AllWithOpts(ctx, hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: impl.RoleLabel + "=" + role + "," + impl.ClusterLabel + "=" + clusterName,
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
