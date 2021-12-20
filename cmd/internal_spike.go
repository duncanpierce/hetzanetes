package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

// Temporary command to explore communication with the Kubernetes API from within the cluster
// This will be removed in future

func Spike() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "spike",
		Short:   "Temporary in-cluster exploratory tool",
		Long:    "Temporary in-cluster exploratory tool",
		Example: "  hetzanetes spike",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			config, err := rest.InClusterConfig()
			if err != nil {
				return err
			}
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}

			fmt.Printf("Starting watch\n")
			nodeWatch, err := clientset.CoreV1().Nodes().Watch(ctx, metav1.ListOptions{})
			if err != nil {
				return err
			}
			// TODO register a termination signal handler and listen to it here
			for event := range nodeWatch.ResultChan() {
				node, ok := event.Object.(*v1.Node)
				if ok {
					fmt.Printf("Node %s %s at %s\n", node.Name, event.Type, time.Now().Format("2006-01-02 15:04:05"))
					fmt.Printf("  Phase: %s, cloud provider: %s, unschedulable: %t\n", node.Status.Phase, node.Spec.ProviderID, node.Spec.Unschedulable)
					fmt.Printf("  Addresses:\n")
					for _, address := range node.Status.Addresses {
						fmt.Printf("    %s: %s\n", address.Type, address.Address)
					}
					fmt.Printf("  Annotations:\n")
					for k, v := range node.Annotations {
						fmt.Printf("    %s: %s\n", k, v)
					}
					fmt.Printf("  Capacity:\n")
					for k, q := range node.Status.Capacity {
						fmt.Printf("    %s: %s\n", k, q.String())
					}
					fmt.Printf("  Allocatable:\n")
					for k, q := range node.Status.Allocatable {
						fmt.Printf("    %s: %s\n", k, q.String())
					}
				}
			}
			return nil
		},
	}
	return cmd
}
