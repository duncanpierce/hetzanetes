package cmd

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/login"
	"github.com/spf13/cobra"
)

func Ssh() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ssh",
		Short: "Temporary in-cluster exploratory tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			public, private, err := login.NewSshKey()
			if err != nil {
				return err
			}
			fmt.Printf("%s\n\n\n%s\n", public, private)
			return nil
		},
	}
	return cmd
}
