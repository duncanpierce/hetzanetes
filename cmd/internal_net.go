package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"net"
)

func Net() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "net",
		Short: "List network interfaces",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			interfaces, err := net.Interfaces()
			if err != nil {
				return err
			}
			for _, i := range interfaces {
				addrs, err := i.Addrs()
				if err != nil {
					return err
				}
				for _, addr := range addrs {
					ip, net, err := net.ParseCIDR(addr.String())
					if err != nil {
						return err
					}
					log.Printf("%s: %s %s %s\n", i.Name, addr.String(), ip.String(), net.String())
				}
			}
			return nil
		},
	}

	return cmd
}
