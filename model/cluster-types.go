package model

import (
	"github.com/Masterminds/semver"
	"github.com/duncanpierce/hetzanetes/label"
)

type (
	Clusters []*Cluster

	Cluster struct {
		ApiVersion string         `json:"apiVersion,omitempty"`
		Kind       string         `json:"kind,omitempty"`
		Metadata   *Metadata      `json:"metadata,omitempty"`
		Spec       *Spec          `json:"spec,omitempty"`
		Status     *ClusterStatus `json:"status,omitempty"`
	}

	Metadata struct {
		Name string `json:"name,omitempty"`
	}

	Actions interface {
		GetServer(name string, apiServer bool, kubernetesVersion *semver.Version) (*NodeStatus, error)
		GetReleaseChannels() (ReleaseChannelStatuses, error)
		CreateServer(name string, serverType string, image string, location string, privateNetworkId string, firewallIds []string, labels label.Labels, sshKeys []string, cloudInit string) (cloudId string, err error)
		DeleteServer(cloudId string) (notFound bool)
		DrainNode(node NodeStatus) error
		CheckNodeReady(node NodeStatus) bool
		CheckNoNode(name string) bool
		DeleteNode(node NodeStatus) error
		SaveStatus(clusterName string, clusterStatus *ClusterStatus) error
	}
)
