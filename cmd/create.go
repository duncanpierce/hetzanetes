package cmd

import (
	"errors"
	"fmt"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/hcloud_client"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/model"
	"github.com/duncanpierce/hetzanetes/tmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/yaml"
	"net"
	"os"
	"strconv"
)

func Create() *cobra.Command {
	var clusterYamlFilename string

	ipRange := net.IPNet{IP: net.IP{10, 0, 0, 0}, Mask: net.IPMask{255, 255, 0, 0}}
	sshPort := "22"
	k3sApiPort := "6443"

	cmd := &cobra.Command{
		Use:              "create [CLUSTER-NAME] [FLAGS]",
		Short:            "Create a new cluster",
		Long:             "Create a new Hetzanetes cluster in a new private network.",
		Example:          "  hetzanetes create [CLUSTER-NAME]",
		TraverseChildren: true,
		Args:             cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var cluster model.Cluster
			var clusterYaml []byte
			var err error

			if clusterYamlFilename != "" {
				if len(args) != 0 {
					return errors.New("cluster name must not be provided when using a cluster YAML file")
				}
				clusterYaml, err = os.ReadFile(clusterYamlFilename)
				if err != nil {
					return err
				}
			} else {
				if len(args) < 1 {
					return errors.New("must provide a cluster name")
				}
				clusterYaml, err = tmpl.DefaultClusterFile(args[0])
				if err != nil {
					return err
				}
			}
			err = yaml.Unmarshal(clusterYaml, &cluster)
			if err != nil {
				return err
			}

			firstApiServerNodeSet := cluster.FirstApiServerNodeSet()
			if firstApiServerNodeSet == nil {
				return errors.New("cluster specifies no API servers")
			}

			c := hcloud_client.New()
			apiToken := env.HCloudToken()
			labels := label.Labels{}
			labels[label.ClusterNameLabel] = cluster.Metadata.Name

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
				Name:    cluster.Metadata.Name,
				IPRange: &ipRange,
				Subnets: subnets,
				Routes:  nil,
				Labels:  networkLabels,
			})
			if err != nil {
				return err
			}
			fmt.Printf("Created network %d %s (%s)\n", network.ID, network.Name, network.IPRange.String())

			serverConfig := tmpl.ClusterConfig{
				HetznerApiToken:   apiToken,
				ClusterName:       cluster.Metadata.Name,
				ClusterNetworkId:  strconv.Itoa(network.ID),
				PrivateIpRange:    ipRange.String(),
				PodIpRange:        "10.42.0.0/16",
				ServiceIpRange:    "10.43.0.0/16",
				InstallDirectory:  "/var/opt/hetzanetes",
				K3sReleaseChannel: cluster.Spec.Versions.GetKubernetes(),
				HetzanetesTag:     cluster.Spec.Versions.GetHetzanetes(),
				ClusterYaml:       string(clusterYaml),
			}
			cloudInit := tmpl.Cloudinit(serverConfig, "create.yaml")

			serverType, _, err := c.ServerType.GetByName(c, firstApiServerNodeSet.ServerType)
			if err != nil {
				return err
			}
			image, _, err := c.Image.GetByName(c, cluster.Spec.Versions.GetBaseImage())
			if err != nil {
				return err
			}

			_, allIPv4, _ := net.ParseCIDR("0.0.0.0/0")
			_, allIPv6, _ := net.ParseCIDR("::/0")
			clusterSelector := label.ClusterNameLabel + "==" + cluster.Metadata.Name
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
				Name:  cluster.Metadata.Name + "-worker",
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
				Name:  cluster.Metadata.Name + "-api",
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

			sshKeys, err := c.SSHKey.All(c)
			if err != nil {
				return err
			}

			t := true
			serverCreateResult, _, err := c.Server.Create(c, hcloud.ServerCreateOpts{
				Name:             firstApiServerNodeSet.ServerName(cluster.Metadata.Name, 1),
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
	cmd.Flags().StringVarP(&clusterYamlFilename, "filename", "f", "", "Name of YAML file specifying cluster configuration")
	return cmd
}
