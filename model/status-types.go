package model

import (
	"github.com/Masterminds/semver"
	"github.com/duncanpierce/hetzanetes/label"
	"time"
)

type (
	ClusterStatus struct {
		VersionStatus   `json:"versions"`
		ClusterNetwork  ClusterNetworkStatus `json:"clusterNetwork"`
		BaseImage       string               `json:"baseImage,omitempty"`
		NodeSetStatuses `json:"nodeSets,omitempty"`
	}
	ClusterNetworkStatus struct {
		CloudId string `json:"cloudId,omitempty"`
		IpRange string `json:"ipRange"`
	}
	NodeSetStatuses []*NodeSetStatus
	NodeSetStatus   struct {
		Name         string `json:"name"`
		Generation   int    `json:"generation"`
		NodeStatuses `json:"nodes,omitempty"`
	}
	NodeStatuses   []NodeStatus
	NodeStatusRefs []*NodeStatus
	NodeStatus     struct {
		Name         string          `json:"name"`
		ServerType   string          `json:"serverType"`
		Location     string          `json:"location"`
		Created      time.Time       `json:"created"`
		CloudId      string          `json:"cloudId,omitempty"`
		ClusterIP    string          `json:"clusterIP,omitempty"`
		BaseImage    string          `json:"baseImage,omitempty"`
		ApiServer    bool            `json:"apiServer,omitempty"`
		Version      *semver.Version `json:"version,omitempty"`
		JoinEndpoint string          `json:"joinEndpoint,omitempty"`
		Phase        `json:"phase"`
		PhaseChanged time.Time `json:"phaseChanged"`
	}
	Phase string

	Actions interface {
		CreateServer(name string, serverType string, image string, location string, privateNetworkId string, firewallIds []string, labels label.Labels, sshKeys []string, cloudInit string) (cloudId string, err error)
		DeleteServer(cloudId string) (notFound bool)
		DrainNode(name string) error
		CheckNodeReady(name string) bool
		CheckNoNode(name string) bool
		DeleteNode(name string)
		SaveStatus(clusterName string, clusterStatus *ClusterStatus) error
	}
)
