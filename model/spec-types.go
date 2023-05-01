package model

import "time"

type (
	Spec struct {
		Versions VersionsSpec `json:"versions,omitempty"`
		NodeSets NodeSetsSpec `json:"nodeSets,omitempty" yaml:"nodeSets"`
	}
	VersionsSpec struct {
		Kubernetes string `json:"kubernetes,omitempty"`
		Hetzanetes string `json:"hetzanetes,omitempty"`
	}
	NodeSetsSpec []*NodeSetSpec
	NodeSetSpec  struct {
		Name        string                     `json:"name"`
		ApiServer   bool                       `json:"apiServer" yaml:"apiServer"`
		Replicas    int                        `json:"replicas"`
		ServerType  string                     `json:"serverType" yaml:"serverType"`
		Locations   []string                   `json:"locations,omitempty"`
		Image       string                     `json:"image"`
		SetUp       []string                   `json:"setUp,omitempty" yaml:"setUp"`
		HealthCheck map[time.Duration][]string `json:"healthCheck,omitempty" yaml:"healthCheck"`
	}
)
