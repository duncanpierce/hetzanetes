package model

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/client/rest"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/login"
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
	sshHostPort := fmt.Sprintf("%s:22", n.ClusterIP)
	var err error

	switch n.Phases.Current().Phase {

	case Create:
		labels := label.Labels{}
		labels.Set(label.Cluster, cluster.Metadata.Name)

		if n.ApiServer {
			labels.Mark(label.ApiServer)
		} else {
			labels.Mark(label.Worker)
		}
		kubernetesVersion := fmt.Sprintf("v%s", cluster.Status.Versions.NewNodeVersion(n.ApiServer).String())
		log.Printf("Version %s chosen for node %s\n", kubernetesVersion, n.Name)

		n.CloudId, n.ClusterIP, err = actions.CreateServer(n.Name, n.ServerType, n.BaseImage, n.Location, cluster.Status.SshPublicKey, cluster.Status.ClusterNetwork.CloudId, nil, labels)
		if err == nil {
			n.SetPhase(Creating, "creating server")
		} else if err == rest.Conflict {
			log.Printf("Conflict: node %s has already been created\n", n.Name)
			existingServers, err := actions.GetServers(cluster.Metadata.Name)
			if err != nil {
				log.Printf("Unable to get existing servers from Hetzner: %s\n", err.Error())
			} else {
				existingServer := existingServers[n.Name]
				n.CloudId = existingServer.CloudId
				n.ClusterIP = existingServer.ClusterIP
				n.SetPhase(Joining, "waiting for previously-created node to join")
			}
		} else {
			log.Printf("error creating server '%s': %s", n.Name, err.Error())
		}

	case Creating:
		if login.PollCloudInit(sshHostPort, env.SshPrivateKey()) {
			n.SetPhase(Install, "server is ready for Kubernetes installation")
		}

	case Install:
		var command login.Command
		if n.ApiServer {
			command = login.AddApiServerCommand(n.JoinEndpoint, env.K3sToken(), n.Version.String())
		} else {
			command = login.AddWorkerCommand(n.JoinEndpoint, env.K3sToken(), n.Version.String())
		}
		err = login.RunCommands(sshHostPort, env.SshPrivateKey(), 3*time.Second, []login.Command{command})
		if err != nil {
			log.Printf("error installing Kubernetes: %s", err.Error())
		} else {
			n.SetPhase(Joining, "waiting for node to join")
		}

	case Joining:
		nodeResource, err := actions.GetNode((*n).Name)
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
		// TODO draining is complete as soon as there are no non-DaemonSet pods - the timeout should be an upper limit (and should be higher)
		if LongerThan(3 * time.Minute)(*n) {
			err := actions.DeleteNode(*n) // TODO might fail if we go straight from Create/Join to Delete with node ever registering
			if err == nil {
				n.SetPhase(Deleting, "")
			} else {
				log.Printf("error deleting node %s\n", n.Name)
			}
		}

	case Deleting:
		_, err = actions.GetNode((*n).Name)
		if err == rest.NotFound && LongerThan(2*time.Minute)(*n) {
			notFound := actions.DeleteServer(*n)
			if notFound {
				n.SetPhase(Deleted, "")
			}
		}

	default:
	}
}
