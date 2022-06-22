package model

import "github.com/Masterminds/semver"

type (
	ReleaseChannelStatuses []*ReleaseChannelStatus

	ReleaseChannelStatus struct {
		Name   string          `json:"name"`
		Latest *semver.Version `json:"latest"`
	}
)

func (r ReleaseChannelStatuses) Named(name string) *ReleaseChannelStatus {
	for _, s := range r {
		if s.Name == name {
			return s
		}
	}
	return nil
}
