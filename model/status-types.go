package model

import (
	"github.com/Masterminds/semver"
	"time"
)

type (
	ClusterStatus struct {
		Versions        VersionStatus        `json:"versions"`
		ClusterNetwork  ClusterNetworkStatus `json:"clusterNetwork"`
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
	NodeStatuses   []*NodeStatus
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
)
