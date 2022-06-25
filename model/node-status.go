package model

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/rest"
	"github.com/duncanpierce/hetzanetes/tmpl"
	"log"
	"time"
)

func (n *NodeStatus) SetPhase(phase Phase, reason string) {
	n.Phases = append(n.Phases, PhaseChange{
		Phase:  phase,
		Reason: reason,
		Time:   time.Now(),
	})
}

func (n *NodeStatus) MakeProgress(cluster *Cluster, actions Actions) {
	var err error

	switch n.Phases.Current().Phase {

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
		kubernetesVersion := fmt.Sprintf("v%s", cluster.Status.Versions.NewNodeVersion(n.ApiServer).String())
		log.Printf("Version %s chosen for node %s\n", kubernetesVersion, n.Name)
		config := tmpl.ClusterConfig{
			KubernetesVersion: kubernetesVersion,
			ApiEndpoint:       n.JoinEndpoint,
			JoinToken:         env.K3sToken(), // TODO this should come from a named Secret
			PrivateIpRange:    cluster.Status.ClusterNetwork.IpRange,
			SshPublicKey:      env.SshPublicKey(), // TODO this should come from a named Secret
		}
		cloudInit := tmpl.Cloudinit(config, templateToUse)
		log.Printf("Cloudinit for new node %s:\n%s\n\n", n.Name, cloudInit)

		sshKeys, err := actions.GetSshKeyIds()
		if err != nil {
			log.Printf("error getting SSH key names: %s\n", err.Error())
		} else {
			n.CloudId, n.ClusterIP, err = actions.CreateServer(n.Name, n.ServerType, n.BaseImage, n.Location, cluster.Status.ClusterNetwork.CloudId, nil, labels, sshKeys, cloudInit)
			if err == nil {
				n.SetPhase(Joining, "waiting for node to join") // TODO once we use SSH, next phase will be Creating
			} else {
				log.Printf("error creating server '%s': %s", n.Name, err.Error())
			}
		}

	case Joining:
		nodeResource, err := actions.GetKubernetesNode(*n)
		if err != nil {
			log.Printf("got error from kubernetes api getting node '%s': %s\n", n.Name, err.Error())
			break
		}
		if nodeResource.IsReady() {
			n.SetPhase(Active, "node has joined")
		}

	case Delete:
		err = actions.DrainNode(*n) // TODO might fail if we go straight from Create/Join to Delete with node ever registering - even if we check whether node has registered and answer is no, we still can't proceed because it's racing us
		if err == nil {
			n.SetPhase(Draining, "")
		} else {
			log.Printf("error draining node '%s': %s", n.Name, err.Error())
		}

	case Draining:
		if LongerThan(3 * time.Minute)(*n) {
			err := actions.DeleteNode(*n) // TODO might fail if we go straight from Create/Join to Delete with node ever registering
			if err == nil {
				n.SetPhase(Deleting, "")
			} else {
				log.Printf("error deleting node %s\n", n.Name)
			}
		}

	case Deleting:
		_, err = actions.GetKubernetesNode(*n)
		if err == rest.NotFound {
			notFound := actions.DeleteServer(*n)
			if notFound {
				n.SetPhase(Deleted, "")
			}
		}

	default:
	}
}
