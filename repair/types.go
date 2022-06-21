package repair

import (
	"github.com/Masterminds/semver"
	"time"
)

type (
	ClusterStatus struct {
		VersionStatus         `json:"versionStatus"`
		ClusterNetworkCloudId string `json:"clusterNetworkCloudId,omitempty"`
		CloudBaseImage        string `json:"cloudBaseImage,omitempty"`
		NodeSetStatuses       `json:"nodeSets,omitempty"`
	}
	NodeSetStatuses []NodeSetStatus
	NodeSetStatus   struct {
		Name             string   `json:"name"`
		Generation       int      `json:"generation"`
		FirewallCloudIds []string `json:"firewallCloudIds,omitempty"`
		NodeStatuses     `json:"nodes,omitempty"`
	}
	NodeStatuses   []NodeStatus
	NodeStatusRefs []*NodeStatus
	NodeStatus     struct {
		Name              string          `json:"name"`
		ServerType        string          `json:"serverType"`
		Location          string          `json:"location"`
		Created           time.Time       `json:"created"`
		CloudId           string          `json:"cloudId,omitempty"`
		BaseImage         string          `json:"baseImage,omitempty"`
		ApiServer         bool            `json:"apiServer,omitempty"`
		KubernetesVersion *semver.Version `json:"kubernetesVersion,omitempty"`
		ClusterIP         string          `json:"clusterIP,omitempty"` // from https://api.hetzner.cloud/v1/servers/{id} .server.private_net[].ip where .network matches cluster network id
		Phase             `json:"phase"`
		PhaseChanged      time.Time `json:"phaseChanged"`
	}
	Phase string

	Actions interface {
		CreateServer(name string, serverType string, location string, privateNetworkId string, firewallIds []string) (cloudId string, clusterIP string, err error)
		DrainNode(name string) error
		DeleteServer(cloudId string) (notFound bool)
		CheckNodeReady(name string) bool
		CheckNoNode(name string) bool
		DeleteNode(name string)
		// etc
	}
)
