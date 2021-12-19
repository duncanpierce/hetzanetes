package cmd

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/client"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"time"
)

func Watch(c client.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch actions in Hetzner API",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ticker := time.NewTicker(1 * time.Second)
			for {
				actions, _, err := c.Action.List(c, hcloud.ActionListOpts{
					Status: []hcloud.ActionStatus{hcloud.ActionStatusRunning},
				})
				fmt.Printf("------------------\n")
				if err != nil {
					fmt.Printf("error: %s\n", err.Error())
				}
				for _, action := range actions {
					fmt.Printf("action: %#v\n\n", action)
				}
				select {
				case <-ticker.C:
					break
				case <-c.Done():
					return c.Err()
				}
			}

			return nil
		},
	}

	return cmd
}
