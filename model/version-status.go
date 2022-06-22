package model

import (
	"fmt"
	"github.com/Masterminds/semver"
)

type (
	VersionStatus struct {
		TargetVersion  *semver.Version            `json:"targetVersion,omitempty"`
		NodeVersions   VersionRange               `json:"nodeVersions"`
		ApiVersions    VersionRange               `json:"apiVersions"`
		WorkerVersions VersionRange               `json:"workerVersions"`
		Channels       map[string]*semver.Version `json:"channels,omitempty"` // gathered from https://update.k3s.io/v1-release/channels
	}

	VersionRange struct {
		Min *semver.Version `json:"min,omitempty"`
		Max *semver.Version `json:"max,omitempty"`
	}
)

func VersionMin(a, b *semver.Version) *semver.Version {
	if a != nil && a.LessThan(b) {
		return a
	}
	return b
}

func VersionMax(a, b *semver.Version) *semver.Version {
	if a != nil && a.GreaterThan(b) {
		return a
	}
	return b
}

func (v VersionRange) Same() bool {
	return v.Max.Equal(v.Min)
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
	return v.ApiVersions.Min
}

func (v VersionStatus) NewApiNodeVersion() *semver.Version {
	// If there are a mix of different versions, drive new API servers towards the same version
	if !v.ApiVersions.Same() {
		return v.ApiVersions.Max
	}

	// Target version can't be a downgrade for any node in the cluster
	if v.TargetVersion.Major() < v.NodeVersions.Max.Major() || v.TargetVersion.Minor() <= v.NodeVersions.Max.Minor() {
		return v.NodeVersions.Max
	}

	// Treat the max allowable version as a minor increment above the lowest version of any node in the cluster
	maxAllowable := VersionMin(v.ApiVersions.Min, v.WorkerVersions.Min).IncMinor()

	// If target version satisfies the max allowed, use that
	if v.TargetVersion.Major() == maxAllowable.Major() && v.TargetVersion.Minor() <= maxAllowable.Minor() {
		return v.TargetVersion
	}

	// Otherwise, we can't directly upgrade because there is too much version skew, so look for the current release in the channel of the max allowed version
	upgradeVersion, found := v.Channels[fmt.Sprint("v%d.%d", maxAllowable.Major(), maxAllowable.Minor())]
	if found {
		return upgradeVersion
	}

	// If all that fails, there is nothing we can do other than try to match the max API server version
	return v.ApiVersions.Max
}
