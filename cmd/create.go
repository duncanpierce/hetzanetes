package cmd

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/hcloud_client"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/tmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"net"
)

// TODO add options to --protected, --backups to enable protection and backups
// TODO maybe protected should be the default
func Create() *cobra.Command {
	var dryRun bool
	var clusterName string
	var ipRange net.IPNet
	var labelsMap map[string]string
	var serverType string
	var osImage string
	var k3sReleaseChannel string

	cmd := &cobra.Command{
		Use:              "create [FLAGS]",
		Short:            "Create a new cluster",
		Long:             "Create a new Hetzanetes cluster in a new private network.",
		Example:          "  hetzanetes create --name=cluster-1",
		TraverseChildren: true,
		Args:             cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c := hcloud_client.New()
			apiToken := env.HCloudToken()
			var labels label.Labels = labelsMap
			labels[label.ClusterNameLabel] = clusterName

			serverConfig := tmpl.ClusterConfig{
				HetznerApiToken:   apiToken,
				ClusterName:       clusterName,
				PrivateIpRange:    ipRange.String(),
				PodIpRange:        "10.42.0.0/16",
				ServiceIpRange:    "10.43.0.0/16",
				InstallDirectory:  "/var/opt/hetzanetes",
				ServerType:        serverType,
				K3sReleaseChannel: k3sReleaseChannel,
			}
			cloudInit := tmpl.Cloudinit(serverConfig, "create.yaml")

			if dryRun {
				fmt.Printf("Would create server with cloud-init file\n%s\n", cloudInit)
				return nil
			}

			// TODO check for name collisions on network and API server before starting, and also on server and network labels
			// TODO split this out behind a driver interface to allow --dry-run

			subnets := []hcloud.NetworkSubnet{
				{
					Type:        hcloud.NetworkSubnetTypeCloud,
					IPRange:     &ipRange,
					NetworkZone: hcloud.NetworkZoneEUCentral,
					Gateway:     nil,
				},
			}

			// TODO protect this network - it could be difficult to repair if deleted (e.g. server gets a new interface flannel doesn't know about)
			networkLabels := labels.Copy().Mark(label.PrivateNetworkLabel)
			network, _, err := c.Network.Create(c, hcloud.NetworkCreateOpts{
				Name:    clusterName,
				IPRange: &ipRange,
				Subnets: subnets,
				Routes:  nil,
				Labels:  networkLabels,
			})
			if err != nil {
				return err
			}
			fmt.Printf("Created network %s (%s)\n", network.Name, network.IPRange.String())

			serverType, _, err := c.ServerType.GetByName(c, serverType)
			if err != nil {
				return err
			}
			image, _, err := c.Image.GetByName(c, osImage)
			if err != nil {
				return err
			}

			_, allIPv4, _ := net.ParseCIDR("0.0.0.0/0")
			_, allIPv6, _ := net.ParseCIDR("::/0")
			sshPort := "22"
			k3sApiPort := "6443"
			clusterSelector := label.ClusterNameLabel + "==" + clusterName
			firewallRules := []hcloud.FirewallRule{
				{
					Protocol:  hcloud.FirewallRuleProtocolICMP,
					SourceIPs: []net.IPNet{*allIPv4, *allIPv6},
					Direction: hcloud.FirewallRuleDirectionIn,
				},
				{
					Port:      &sshPort,
					Protocol:  hcloud.FirewallRuleProtocolTCP,
					SourceIPs: []net.IPNet{*allIPv4, *allIPv6},
					Direction: hcloud.FirewallRuleDirectionIn,
				},
			}
			_, _, err = c.Firewall.Create(c, hcloud.FirewallCreateOpts{
				Name:  clusterName + "-worker",
				Rules: firewallRules,
				ApplyTo: []hcloud.FirewallResource{
					{
						Type: hcloud.FirewallResourceTypeLabelSelector,
						LabelSelector: &hcloud.FirewallResourceLabelSelector{
							Selector: clusterSelector + "," + label.WorkerLabel,
						},
					},
				},
			})
			if err != nil {
				return err
			}
			// TODO API port not required if we use a load balancer and access from private IP
			firewallRules = append(firewallRules, hcloud.FirewallRule{
				Port:      &k3sApiPort,
				Protocol:  hcloud.FirewallRuleProtocolTCP,
				SourceIPs: []net.IPNet{*allIPv4, *allIPv6},
				Direction: hcloud.FirewallRuleDirectionIn,
			})
			_, _, err = c.Firewall.Create(c, hcloud.FirewallCreateOpts{
				Name:  clusterName + "-api",
				Rules: firewallRules,
				ApplyTo: []hcloud.FirewallResource{
					{
						Type: hcloud.FirewallResourceTypeLabelSelector,
						LabelSelector: &hcloud.FirewallResourceLabelSelector{
							Selector: clusterSelector + "," + label.ApiServerLabel,
						},
					},
				},
			})
			if err != nil {
				return err
			}

			// TODO allow a label selector to select keys to use (repair will keep it up to date)
			sshKeys, err := c.SSHKey.All(c)
			if err != nil {
				return err
			}

			// Hetzner recommend specifying locations rather than datacenters: https://docs.hetzner.cloud/#servers-create-a-server
			// TODO add --regions option
			t := true
			serverCreateResult, _, err := c.Server.Create(c, hcloud.ServerCreateOpts{
				Name:             clusterName + "-api-1",
				ServerType:       serverType,
				Image:            image,
				SSHKeys:          sshKeys,
				Location:         nil,
				UserData:         cloudInit,
				StartAfterCreate: &t,
				Labels:           labels.Copy().Mark(label.ApiServerLabel), // TODO --segregate-api to remove this and taint the api server (or have repair do it)
				Networks:         []*hcloud.Network{network},
			})
			if err != nil {
				return err
			}
			fmt.Printf("Creating server %s in %s...\n", serverCreateResult.Server.Name, serverCreateResult.Server.Datacenter.Name)
			return c.Await(serverCreateResult.Action)
		},
	}
	cmd.Flags().StringVar(&clusterName, "name", "", "Cluster name (required)")
	cmd.MarkFlagRequired("name")
	cmd.Flags().IPNetVar(&ipRange, "cluster-ip-range", net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.IPMask{255, 255, 0, 0}}, "Cluster network IP range")
	// TODO remove cluster-ip-range option? make it an attribute of the network provider?
	// TODO allow create-time-only configuration of pod and service IP ranges? might be easier to leave it on defaults
	cmd.Flags().StringToStringVar(&labelsMap, "label", map[string]string{}, "User-defined labels ('key=value') (can be specified multiple times)")
	cmd.Flags().StringVar(&serverType, "server-type", "cx11", "Server type")
	cmd.Flags().StringVar(&osImage, "os-image", "ubuntu-20.04", "Operating system image")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show what would be done without taking any action")
	cmd.Flags().StringVar(&k3sReleaseChannel, "channel", "stable", "K3s release channel")

	return cmd
}
