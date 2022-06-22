package model

import (
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/tmpl"
	"time"
)

func (n *NodeStatus) SetPhase(phase Phase) {
	n.Phase = phase
	n.PhaseChanged = time.Now()
}

func (n *NodeStatus) MakeProgress(cluster *Cluster, actions Actions) {
	var err error

	switch n.Phase {

	case Create:
		var templateToUse string
		labels := label.Labels{}
		labels.Set(label.ClusterNameLabel, cluster.Metadata.Name)

		if n.ApiServer {
			templateToUse = "add-api-server.yaml"
			labels.Mark(label.ApiServerLabel)
		} else {
			templateToUse = "add-worker.yaml"
			labels.Mark(label.WorkerLabel)
		}
		// TODO these cloudinit scripts need .ApiEndpoint, .JoinToken, .PrivateIpRange, .K3sReleaseChannel <- should be VERSION
		config := tmpl.ClusterConfig{
			KubernetesVersion: cluster.Status.TargetVersion.String(),
			ApiEndpoint:       n.JoinEndpoint,
			JoinToken:         env.K3sToken(), // TODO this should come from a named Secret
			PrivateIpRange:    cluster.Status.ClusterNetwork.IpRange,
		}
		cloudInit := tmpl.Cloudinit(config, templateToUse)

		n.CloudId, err = actions.CreateServer(n.Name, n.ServerType, n.BaseImage, n.Location, cluster.Status.ClusterNetwork.CloudId, nil, labels, nil, cloudInit)
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
