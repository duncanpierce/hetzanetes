package model

type (
	ClusterList struct {
		Items Clusters `json:"items"`
	}

	NodeResource struct {
		Status NodeResourceStatus `json:"status,omitempty"`
	}

	NodeResourceStatus struct {
		Conditions []NodeResourceCondition `json:"conditions,omitempty"`
	}

	NodeResourceCondition struct {
		Status string `json:"status,omitempty"`
		Type   string `json:"type,omitempty"`
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
