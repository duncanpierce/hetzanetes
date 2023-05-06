package model

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/duncanpierce/hetzanetes/model/k3s"
)

type (
	VersionStatus struct {
		Target   *semver.Version            `json:"target,omitempty"`
		Nodes    VersionRange               `json:"nodes"`
		Api      VersionRange               `json:"api"`
		Workers  VersionRange               `json:"workers"`
		Channels k3s.ReleaseChannelStatuses `json:"channels,omitempty"`
	}

	VersionRange struct {
		Min *semver.Version `json:"min,omitempty"`
		Max *semver.Version `json:"max,omitempty"`
	}
)

func VersionMin(a, b *semver.Version) *semver.Version {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.LessThan(b) {
		return a
	}
	return b
}

func VersionMax(a, b *semver.Version) *semver.Version {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.GreaterThan(b) {
		return a
	}
	return b
}

func (v VersionRange) Same() bool {
	if v.Max == nil {
		return v.Min == nil
	} else if v.Min == nil {
		return false
	} else {
		return v.Max.Equal(v.Min)
	}
}

func (v VersionRange) MergeRange(other VersionRange) VersionRange {
	return VersionRange{
		Min: VersionMin(v.Min, other.Min),
		Max: VersionMax(v.Max, other.Max),
	}
}

func (v VersionRange) MergeVersion(other *semver.Version) VersionRange {
	return v.MergeRange(VersionRange{other, other})
}

func (v VersionStatus) NewNodeVersion(apiServer bool) *semver.Version {
	if apiServer {
		return v.NewApiNodeVersion()
	}
	return v.NewWorkerNodeVersion()
}

func (v VersionStatus) NewWorkerNodeVersion() *semver.Version {
	// Worker node can never have a higher version than any API node
	return v.Api.Min
}

func (v VersionStatus) NewApiNodeVersion() *semver.Version {
	// If there is a mix of different versions, drive new API servers towards the same version
	if !v.Api.Same() {
		return v.Api.Max
	}

	// Target major.minor version can't be a downgrade for any node in the cluster (but patch level can be downgraded)
	if v.Target.Major() < v.Nodes.Max.Major() || v.Target.Minor() <= v.Nodes.Max.Minor() {
		return v.Nodes.Max
	}

	// The max allowable version is a minor increment above the lowest version of any node in the cluster
	maxAllowable := v.Nodes.Min.IncMinor()

	// If target version satisfies the max allowed, use that
	if v.Target.Major() == maxAllowable.Major() && v.Target.Minor() <= maxAllowable.Minor() {
		return v.Target
	}

	// Otherwise, we can't directly upgrade to target version because there is too much version skew, so look for the current release in the channel of the max version that is allowed
	channel := v.Channels.Named(fmt.Sprintf("v%d.%d", maxAllowable.Major(), maxAllowable.Minor()))
	if channel != nil {
		return channel.Latest
	}

	// If all that fails, there is nothing we can do other than try to match the max API server version
	return v.Api.Max
}

func (v *VersionStatus) UpdateReleaseChannels(releaseChannel string, actions Actions) error {
	channels, err := actions.GetReleaseChannels()
	if err != nil {
		return err
	}
	v.Channels = channels

	targetVersion := channels.Named(releaseChannel)
	if targetVersion == nil {
		return fmt.Errorf("cannot find a release for channel %s", releaseChannel)
	}
	v.Target = targetVersion.Latest
	return nil
}
