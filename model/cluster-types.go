package model

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
)
