package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Temporary command to explore communication with the Kubernetes API from within the cluster
// This will be removed in future

func Spike(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "spike",
		Short:   "Temporary in-cluster exploratory tool",
		Long:    "Temporary in-cluster exploratory tool",
		Example: "  hetzanetes spike",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := rest.InClusterConfig()
			if err != nil {
				return err
			}
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}
			nodeList, err := clientset.CoreV1().Nodes().List(ctx, v1.ListOptions{})
			if err != nil {
				return err
			}
			for _, node := range nodeList.Items {
				fmt.Printf("Node %s\n", node.ObjectMeta.Name)
				for _, address := range node.Status.Addresses {
					fmt.Printf("  %s: %s\n", address.Type, address.Address)
				}
			}
			return nil
		},
	}
	return cmd
}
