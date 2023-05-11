package cmd

import (
	"errors"
	"fmt"
	hcloudClient "github.com/duncanpierce/hetzanetes/client/hcloud"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/login"
	"github.com/duncanpierce/hetzanetes/model"
	"github.com/duncanpierce/hetzanetes/tmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

func Create() *cobra.Command {
	var clusterYamlFilename string
	var installHetzanetesVersion string

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
				clusterYaml = []byte(fmt.Sprintf(tmpl.DefaultClusterYaml, args[0]))
				if err != nil {
					return err
				}
			}
			err = yaml.Unmarshal(clusterYaml, &cluster)
			if err != nil {
				return err
			}

			bootstrapServerName, err := cluster.BootstrapApiServerName()
			if err != nil {
				return err
			}

			c := hcloudClient.New()
			labels := label.Labels{}
			labels[label.Cluster] = cluster.Metadata.Name

			// TODO check for name collisions on network and API server before starting, and also on server and network labels

			subnets := []hcloud.NetworkSubnet{
				{
					Type:        hcloud.NetworkSubnetTypeCloud,
					IPRange:     &ipRange,
					NetworkZone: hcloud.NetworkZoneEUCentral,
					Gateway:     nil,
				},
			}

			// TODO protect this network - it could be difficult to repair if deleted (e.g. server gets a new interface flannel doesn't know about)
			networkLabels := labels.Copy().Mark(label.PrivateNetwork)
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
			log.Printf("Created network %d %s (%s)\n", network.ID, network.Name, network.IPRange.String())

			sshPublicKey, sshPrivateKey, err := login.NewSshKey()
			if err != nil {
				return err
			}

			_, allIPv4, _ := net.ParseCIDR("0.0.0.0/0")
			_, allIPv6, _ := net.ParseCIDR("::/0")
			clusterSelector := label.Cluster + "==" + cluster.Metadata.Name
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
							Selector: clusterSelector + "," + label.Worker,
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
							Selector: clusterSelector + "," + label.ApiServer,
						},
					},
				},
			})
			if err != nil {
				return err
			}

			actions := model.NewClusterActions()
			err = cluster.Create(actions, strconv.Itoa(network.ID), sshPublicKey)
			if err != nil {
				return err
			}
			nodeStatuses, err := actions.GetServers(cluster.Metadata.Name)
			if err != nil {
				return err
			}
			bootstrapNodeStatus := nodeStatuses[bootstrapServerName]
			sshHostPort := fmt.Sprintf("%s:22", bootstrapNodeStatus.PublicIPv4)
			login.AwaitCloudInit(sshHostPort, sshPrivateKey)
			log.Printf("bootstrapping cluster\n")
			commands := login.CreateClusterCommands(clusterYaml, cluster.Metadata.Name, strconv.Itoa(network.ID), ipRange.String(), cluster.Spec.Versions.GetKubernetes(), installHetzanetesVersion, env.HCloudToken(), sshPrivateKey, sshPublicKey)
			return login.RunCommands(sshHostPort, sshPrivateKey, 3*time.Second, commands)
		},
	}
	cmd.Flags().StringVarP(&clusterYamlFilename, "filename", "f", "", "Name of YAML file specifying cluster configuration")
	cmd.Flags().StringVarP(&installHetzanetesVersion, "hetzanetes-version", "v", "latest", "Version of Hetzanetes to install in the cluster ('none' to skip install for testing)")
	return cmd
}
