package model

import (
	"github.com/Masterminds/semver"
	"github.com/duncanpierce/hetzanetes/label"
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
		GetBootstrapServer(name string, apiServer bool, kubernetesVersion *semver.Version) (*NodeStatus, error)
		GetReleaseChannels() (ReleaseChannelStatuses, error)
		CreateServer(name string, serverType string, image string, location string, privateNetworkId string, firewallIds []string, labels label.Labels, sshKeyIds []int, cloudInit string) (cloudId string, clusterIP string, err error)
		DeleteServer(node NodeStatus) (notFound bool)
		DrainNode(node NodeStatus) error
		GetKubernetesNode(node NodeStatus) (*NodeResource, error)
		DeleteNode(node NodeStatus) error
		SaveStatus(clusterName string, clusterStatus *ClusterStatus) error
		GetSshKeyIds() ([]int, error)
	}
)
