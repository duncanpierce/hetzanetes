package repair

import (
	"sort"
	"time"
)

type Filter func(n NodeStatus) bool

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

func (n NodeStatusRefs) GetVersionRange() (v VersionRange) {
	for _, node := range n {
		v = v.MergeVersion(node.KubernetesVersion)
	}
	return
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
