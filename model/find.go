package model

import (
	"sort"
	"time"
)

type Filter func(n NodeStatus) bool

func (c *ClusterStatus) Find(filters ...Filter) (nodes NodeStatusRefs) {
	for _, nodeSet := range c.NodeSetStatuses {
		nodes = append(nodes, nodeSet.Find(filters...)...)
	}
	return
}

func (n *NodeStatuses) Find(filters ...Filter) (nodes NodeStatusRefs) {
	for _, node := range *n {
		for _, filter := range filters {
			if !filter(node) {
				break
			}
			nodes = append(nodes, &node)
		}
	}
	return
}

func (n NodeStatusRefs) SortByPhase() {
	sort.SliceStable(n, func(i, j int) bool {
		return n[i].Phase.Compare(n[j].Phase) < 0
	})
}

// SortByRecency returns most recent first
func (n NodeStatusRefs) SortByRecency() {
	sort.SliceStable(n, func(i, j int) bool {
		return (n[j].PhaseChanged).Before(n[i].PhaseChanged)
	})
}

func (n NodeStatusRefs) SetPhase(phase Phase) {
	for i := 0; i < len(n); i++ {
		n[i].SetPhase(phase)
	}
}

func (n NodeStatusRefs) MakeProgress(cluster *Cluster, actions Actions) {
	for _, node := range n {
		node.MakeProgress(cluster, actions)
	}
}

func (n NodeStatusRefs) GetVersionRange() (v VersionRange) {
	for _, node := range n {
		v = v.MergeVersion(node.Version)
	}
	return
}

func MatchAll() Filter {
	return func(node NodeStatus) bool {
		return true
	}
}

func InPhase(phases ...Phase) Filter {
	return func(node NodeStatus) bool {
		for _, phase := range phases {
			if node.Phase == phase {
				return true
			}
		}
		return false
	}
}

func PhaseUpTo(phase Phase) Filter {
	return func(node NodeStatus) bool {
		return node.Phase.Compare(phase) <= 0
	}
}

func LongerThan(d time.Duration) Filter {
	return func(node NodeStatus) bool {
		return node.PhaseChanged.Before(time.Now().Add(-d))
	}
}

func IsApiServer(t bool) Filter {
	return func(node NodeStatus) bool {
		return node.ApiServer == t
	}
}
