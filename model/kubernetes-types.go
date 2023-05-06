package model

import (
	"github.com/Masterminds/semver"
)

type (
	ClusterList struct {
		Items Clusters `json:"items"`
	}

	NodeResource struct {
		Spec   *NodeResourceSpec   `json:"spec,omitempty"`
		Status *NodeResourceStatus `json:"status,omitempty"`
	}

	NodeResourceSpec struct {
		Unschedulable bool `json:"unschedulable,omitempty"`
	}

	NodeResourceStatus struct {
		Conditions []NodeResourceCondition `json:"conditions,omitempty"`
		NodeInfo   NodeInfo                `json:"nodeInfo,omitempty"`
	}

	NodeInfo struct {
		Architecture   string          `json:"architecture,omitempty"`
		KubeletVersion *semver.Version `json:"kubeletVersion,omitempty"`
	}

	NodeResourceCondition struct {
		Status string `json:"status,omitempty"`
		Type   string `json:"type,omitempty"`
	}

	PodResourceList struct {
		Items []PodResource `json:"items"`
	}

	PodResource struct {
		Metadata *PodMetadata `json:"metadata,omitempty"`
	}

	PodMetadata struct {
		Name            string               `json:"name,omitempty"`
		Namespace       string               `json:"namespace,omitempty"`
		OwnerReferences []*PodOwnerReference `json:"ownerReferences,omitempty"`
	}

	PodOwnerReference struct {
		ApiVersion string `json:"apiVersion,omitempty"`
		Kind       string `json:"kind,omitempty"`
		Name       string `json:"name,omitempty"`
	}

	PodEviction struct {
		ApiVersion string              `json:"apiVersion"`
		Kind       string              `json:"kind"`
		Metadata   PodEvictionMetadata `json:"metadata,omitempty"`
	}

	PodEvictionMetadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	}
)

func (n NodeResource) IsReady() bool {
	for _, condition := range n.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			return true
		}
	}
	return false
}
