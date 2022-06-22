package model

type (
	Spec struct {
		Versions VersionsSpec `json:"versions,omitempty"`
		NodeSets NodeSetsSpec `json:"nodeSets,omitempty"`
	}
	VersionsSpec struct {
		BaseImage  string `json:"baseImage,omitempty"`
		Kubernetes string `json:"kubernetes,omitempty"`
		Hetzanetes string `json:"hetzanetes,omitempty"`
	}
	NodeSetsSpec []*NodeSetSpec
	NodeSetSpec  struct {
		Name       string   `json:"name"`
		ApiServer  bool     `json:"apiServer"`
		Replicas   int      `json:"replicas"`
		ServerType string   `json:"serverType"`
		Locations  []string `json:"locations,omitempty"`
	}
)

//func (n *NodeSetStatuses) Named(name string) *NodeSetStatus {
//	for _, x := range *n {
//		if x.Name == name {
//			return x
//		}
//	}
//	return nil
//}
