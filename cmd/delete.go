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
			if network == nil {
				return errors.New(fmt.Sprintf("cluster network %s not found", clusterName))
			}
			if network.Labels[impl.RoleLabel] != impl.ClusterRole {
				return errors.New(fmt.Sprintf("the network %s does not have the %s=%s label", clusterName, impl.RoleLabel, impl.ClusterRole))
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

			errs := impl.Errors{}
			for _, server := range servers {
				// TODO retry if deletion fails?
				errs.Catch(client.Server.Delete(ctx, server))
			}
			if len(errs) == 0 {
				errs.Catch(client.Network.Delete(ctx, network))
			}
			return errs
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
