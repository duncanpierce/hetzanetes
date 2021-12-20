package cmd

import (
	"github.com/duncanpierce/hetzanetes/tmpl"
	"github.com/spf13/cobra"
)

// Write files needed to configure the cluster
func Kustomize() *cobra.Command {
	clusterConfig := tmpl.ClusterConfig{
		ClusterName: "",
		PodIpRange:  "",
	}

	cmd := &cobra.Command{
		Use:              "kustomize [FLAGS]",
		Short:            "Write files needed to initialise the cluster",
		TraverseChildren: true,
		Args:             cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			tmpl.WriteKustomizeFiles(clusterConfig)
			return nil
		},
	}

	cmd.Flags().StringVar(&clusterConfig.ClusterName, "name", "", "cluster name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringVar(&clusterConfig.PodIpRange, "pod-ip-range", "", "pod IP range")
	cmd.MarkFlagRequired("pod-ip-range")
	return cmd
}
