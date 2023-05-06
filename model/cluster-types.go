package model

import (
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/model/k3s"
)

type (
	Clusters []*Cluster

	Cluster struct {
		ApiVersion string         `json:"apiVersion,omitempty" yaml:"apiVersion"`
		Kind       string         `json:"kind,omitempty"`
		Metadata   *Metadata      `json:"metadata,omitempty"`
		Spec       *Spec          `json:"spec,omitempty"`
		Status     *ClusterStatus `json:"status,omitempty"`
	}

	Metadata struct {
		Name string `json:"name,omitempty"`
	}

	Actions interface {
		GetServers(clusterName string) (map[string]*NodeStatus, error)
		GetReleaseChannels() (k3s.ReleaseChannelStatuses, error)
		CreateServer(name string, serverType string, image string, location string, sshPublicKey string, privateNetworkId string, firewallIds []string, labels label.Labels) (cloudId string, clusterIP string, err error)
		DeleteServer(node NodeStatus) (notFound bool)
		DrainNode(node NodeStatus) error
		GetNode(name string) (*NodeResource, error)
		DeleteNode(node NodeStatus) error
		SaveStatus(clusterName string, clusterStatus *ClusterStatus) error
		GetSshKeyIds() ([]int, error)
	}
)
