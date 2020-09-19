package main

import (
	"github.com/duncanpierce/hetzanetes/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	var defaultCmd = &cobra.Command{
		Use: "hetzanetes",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	defaultCmd.AddCommand(cmd.List())

	if err := defaultCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
