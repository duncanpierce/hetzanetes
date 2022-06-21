package repair

import "time"

func (n *NodeStatus) SetPhase(phase Phase) {
	n.Phase = phase
	n.PhaseChanged = time.Now()
}

func (n *NodeStatus) NextAction(clusterStatus *ClusterStatus, nodeSetStatus *NodeSetStatus, actions Actions) {
	var err error

	switch n.Phase {

	case Create:
		// TODO pass in cloudinit script
		n.CloudId, n.ClusterIP, err = actions.CreateServer(n.Name, n.ServerType, n.Location, clusterStatus.ClusterNetworkCloudId, nodeSetStatus.FirewallCloudIds)
		if err == nil {
			n.SetPhase(Joining) // TODO once we use SSH, next phase will be Creating
		}

	case Joining:
		ready := actions.CheckNodeReady(n.Name)
		if ready {
			n.SetPhase(Active)
		}

	case Delete:
		err = actions.DrainNode(n.Name) // TODO might fail if we go straight from Create/Join to Delete with node ever registering - even if we check whether node has registered and answer is no, we still can't proceed because it's racing us
		if err == nil {
			n.SetPhase(Draining)
		}

	case Draining:
		if LongerThan(5 * time.Minute)(*n) {
			actions.DeleteNode(n.Name) // TODO might fail if we go straight from Create/Join to Delete with node ever registering
			n.SetPhase(Deleting)
		}

	case Deleting:
		if actions.CheckNoNode(n.Name) {
			notFound := actions.DeleteServer(n.CloudId)
			if notFound {
				n.SetPhase(Deleted)
			}
		}

	default:
	}
}
