package model

const (
	// until we switch to SSH, the flow is simpler and misses some setup steps which are done by Cloudinit

	Create     = Phase("create")     // signal to start provisioning process
	Creating   = Phase("creating")   // Cloud API has been called to create node, attach firewall and network, store SSH key - next action: probe until SSH login works
	Updating   = Phase("updating")   // have logged in, running apt update, apt upgrade, reboot, wait until it disconnects for reboot then go to installing phase
	Installing = Phase("installing") // probe until SSH login works, install K3s using secret, disconnect, wait for Node resource to appear
	Joining    = Phase("joining")    // wait for Node resource to reach Active state (is this step worth having separate from previous? might help with recovering from crash?)
	Active     = Phase("ready")      // node is ready and working - this is the main lifecycle phase
	Replace    = Phase("replace")    // node needs to be replaced
	Delete     = Phase("delete")     // signal to start deletion process
	Draining   = Phase("draining")   // when repair sees unhealthy node, it tells K8s to drain and enters this state
	Deleting   = Phase("deleting")   // when node has been in draining state for 5 minutes, repair tells Hetzner to delete the server
	Deleted    = Phase("deleted")    // server has been deleted from cloud
)

func (p Phase) index() int {
	switch p {
	case Create:
		return 1
	case Creating:
		return 2
	case Updating:
		return 3
	case Installing:
		return 4
	case Joining:
		return 5
	case Active:
		return 6
	case Replace:
		return 7
	case Delete:
		return 8
	case Draining:
		return 9
	case Deleting:
		return 10
	case Deleted:
		return 11
	default:
		return 0
	}
}

func (p Phase) Compare(other Phase) int {
	return other.index() - p.index()
}

func (p PhaseChanges) Current() PhaseChange {
	return p[len(p)-1]
}
