package hcloud_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/tmpl"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"os"
)

type Client struct {
	*hcloud.Client
	context.Context
}

func New() Client {
	return Client{
		Client:  hcloud.NewClient(hcloud.WithToken(env.HCloudToken())),
		Context: context.Background(),
	}
}

func (c Client) Await(action *hcloud.Action) error {
	_, errors := c.Client.Action.WatchProgress(c.Context, action)
	return <-errors
}

func (c Client) CreateServer(apiToken string, clusterName string, nodeSetName string, apiServer bool, serverTypeName string, osImageName string, generationNumber int, k3sReleaseChannel string) error {
	labels := label.Labels{
		label.ClusterNameLabel: clusterName,
	}

	network, _, err := c.Network.GetByName(c, clusterName)
	if err != nil {
		return err
	}
	_, labelled := network.Labels[label.PrivateNetworkLabel]
	if !labelled || network.Labels[label.ClusterNameLabel] != clusterName {
		return errors.New(fmt.Sprintf("network %s is not labelled as a cluster", network.Name))
	}

	nodeName := fmt.Sprintf("%s-%s-%d", clusterName, nodeSetName, generationNumber)
	ipRange := network.Subnets[0].IPRange

	clusterConfig := tmpl.ClusterConfig{
		JoinToken:         os.Getenv("K3S_TOKEN"),
		ApiEndpoint:       os.Getenv("K3S_URL"),
		HetznerApiToken:   apiToken,
		ClusterName:       clusterName,
		PrivateIpRange:    ipRange.String(),
		K3sReleaseChannel: k3sReleaseChannel,
	}
	templateToUse := "add-worker.yaml"
	if apiServer {
		templateToUse = "add-api-server.yaml"
	}
	cloudInit := tmpl.Cloudinit(clusterConfig, templateToUse)

	serverType, _, err := c.ServerType.GetByName(c, serverTypeName)
	if err != nil {
		return err
	}
	osImage, _, err := c.Image.GetByName(c, osImageName)
	if err != nil {
		return err
	}

	// TODO allow a label selector to select keys to use (repair will keep it up to date)
	sshKeys, err := c.SSHKey.All(c)
	if err != nil {
		return err
	}

	// Hetzner recommend specifying locations rather than datacenters: https://docs.hetzner.cloud/#servers-create-a-server
	t := true
	labelToUse := label.WorkerLabel
	if apiServer {
		labelToUse = label.ApiServerLabel
	}
	server, _, err := c.Server.Create(c, hcloud.ServerCreateOpts{
		Name:             nodeName,
		ServerType:       serverType,
		Image:            osImage,
		SSHKeys:          sshKeys,
		Location:         nil,
		StartAfterCreate: &t,
		UserData:         cloudInit,
		Labels:           labels.Copy().Mark(labelToUse),
		Networks:         []*hcloud.Network{network},
	})
	if err != nil {
		return err
	}

	fmt.Printf("Created server %s in %s\n", server.Server.Name, server.Server.Datacenter.Name)
	return nil
}

func (c Client) DrainAndDeleteServer(server *hcloud.Server) error {
	// TODO cordon, drain, pause?, delete node, pause?, delete server
	_, err := c.Server.Delete(c, server)
	return err
}
