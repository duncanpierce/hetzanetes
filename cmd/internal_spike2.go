package cmd

import (
	"github.com/duncanpierce/hetzanetes/k8s_client"
	"github.com/spf13/cobra"
)

// Temporary command to explore communication with the Kubernetes API from within the cluster
// This will be removed in future

func Spike2() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "spike2",
		Short:   "Temporary in-cluster exploratory tool",
		Long:    "Temporary in-cluster exploratory tool",
		Example: "  hetzanetes spike2",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := k8s_client.New()
			_, err := client.GetClusters()
			return err
		},
	}
	return cmd
}
