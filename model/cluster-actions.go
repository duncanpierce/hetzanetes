package model

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/duncanpierce/hetzanetes/catch"
	"github.com/duncanpierce/hetzanetes/client/rest"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/json"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/model/hetzner"
	"github.com/duncanpierce/hetzanetes/model/k3s"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type (
	ClusterActions struct {
		kubernetes *rest.Client
		hetzner    *rest.Client
		k3s        *rest.Client
	}
)

var (
	_              Actions = ClusterActions{}
	strategicMerge         = map[string]string{"Content-Type": "application/merge-patch+json"}
)

func NewClusterActions() *ClusterActions {
	return &ClusterActions{
		kubernetes: NewKubernetes(),
		hetzner:    NewHetzner(env.HCloudToken()),
		k3s:        NewK3s(),
	}
}

func NewKubernetes() *rest.Client {
	cert, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		return nil
	}
	certs := x509.NewCertPool()
	certs.AppendCertsFromPEM(cert)
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            certs,
			},
		},
	}
	tokenFile, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil
	}
	return &rest.Client{
		BaseUrl:     "https://kubernetes.default.svc",
		Http:        client,
		Token:       string(tokenFile),
		LogResponse: true,
	}
}

func NewHetzner(token string) *rest.Client {
	return &rest.Client{
		BaseUrl: "https://api.hetzner.cloud/v1",
		Http:    http.Client{},
		Token:   token,
	}
}

func NewK3s() *rest.Client {
	return &rest.Client{
		BaseUrl: "https://update.k3s.io/v1-release",
		Http:    http.Client{},
	}
}

func (c ClusterActions) GetReleaseChannels() (k3s.ReleaseChannelStatuses, error) {
	response := &k3s.ReleaseChannelsResponse{}
	err := c.k3s.Do(http.MethodGet, "/channels", nil, nil, response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c ClusterActions) GetServers(clusterName string) (map[string]*NodeStatus, error) {
	hetznerServers := &hetzner.ServersResponse{}
	err := c.hetzner.Do(http.MethodGet, fmt.Sprintf("/servers?per_page=50&label_selector=%s==%s", label.Cluster, clusterName), nil, nil, hetznerServers)
	if err != nil {
		return nil, err
	}

	nodeStatuses := map[string]*NodeStatus{}
	for _, server := range hetznerServers.Servers {
		privateNetwork := server.PrivateNets[0]
		nodeStatuses[server.Name] = &NodeStatus{
			Name:       server.Name,
			ServerType: server.ServerType.Name,
			Location:   server.Datacenter.Location.Name,
			CloudId:    strconv.Itoa(server.Id),
			ClusterIP:  privateNetwork.IP,
			PublicIPv4: server.PublicNet.IPv4.IP,
			BaseImage:  server.Image.Name,
			Phases: PhaseChanges{
				PhaseChange{
					Phase:  Create,
					Reason: "existing server",
					Time:   server.Created,
				},
			},
		}
	}

	return nodeStatuses, nil
}

func (c ClusterActions) CreateServer(name string, serverType string, image string, location string, sshPublicKey string, privateNetworkId string, firewallIds []string, labels label.Labels) (cloudId string, clusterIP string, err error) {
	cloudInit := fmt.Sprintf(`#cloud-config

package_update: true
package_upgrade: true
package_reboot_if_required: true
packages:
  - curl
  - jq
users:
  - name: root
    ssh_authorized_keys:
      - %s`, sshPublicKey)

	privateNetworkNumber, _ := strconv.Atoi(privateNetworkId)
	firewalls := []hetzner.FirewallRef{}
	for _, firewallId := range firewallIds {
		firewallNumber, _ := strconv.Atoi(firewallId)
		firewalls = append(firewalls, hetzner.FirewallRef{firewallNumber})
	}

	var sshKeyIds []int
	sshKeyIds, err = c.GetSshKeyIds()
	if err != nil {
		return
	}

	serverRequest := hetzner.CreateServerRequest{
		Name:       name,
		ServerType: serverType,
		Image:      image,
		Location:   location,
		Networks:   []int{privateNetworkNumber},
		Firewalls:  firewalls,
		Labels:     labels,
		SshKeys:    sshKeyIds,
		CloudInit:  cloudInit,
	}
	serverResult := &hetzner.ServerResult{}
	err = c.hetzner.Do(http.MethodPost, "/servers", rest.JSON(), serverRequest, serverResult)
	if err != nil {
		return
	}
	cloudIdNum := serverResult.Server.Id
	log.Printf("Created server %s cloudid=%d\n", name, cloudIdNum)

	// Wait for private network to be attached (among other things)
	err = c.Await("servers", cloudIdNum)
	if err != nil {
		return
	}

	serverResult = &hetzner.ServerResult{}
	err = c.hetzner.Do(http.MethodGet, fmt.Sprintf("/servers/%d", cloudIdNum), nil, nil, serverResult)
	if err != nil {
		return
	}

	return strconv.Itoa(cloudIdNum), serverResult.Server.PrivateNet[0].IP, nil
}

func (c ClusterActions) Await(resourceType string, resourceId int) error {
	for {
		time.Sleep(1 * time.Second)
		result := &hetzner.ActionsResponse{}
		err := c.hetzner.Do(http.MethodGet, fmt.Sprintf("/%s/%d/actions", resourceType, resourceId), nil, nil, result)
		if err != nil {
			log.Printf("error awaiting %s %d: API returned error %s\n", resourceType, resourceId, err.Error())
			return err
		}

		stillRunning := false
		for _, action := range result.Actions {
			if action.Status == "error" {
				return fmt.Errorf("error awaiting %s %d: %s", resourceType, resourceId, action.Error.Message)
			} else if action.Status == "running" {
				stillRunning = true
				break
			}
		}
		if !stillRunning {
			return nil
		}
	}
}

func (c ClusterActions) DeleteServer(node NodeStatus) (notFound bool) {
	if node.CloudId == "" {
		log.Printf("Error: deleting server with no cloudId")
		return false
	} else {
		return c.hetzner.Do(http.MethodDelete, "/servers/"+node.CloudId, nil, nil, nil) == rest.NotFound
	}
}

func (c ClusterActions) DrainNode(node NodeStatus) error {
	log.Printf("Draining node %s\n", json.Format(node))
	nodePatch := &NodeResource{
		Spec: &NodeResourceSpec{
			Unschedulable: true,
		},
	}
	// cordon
	err := c.kubernetes.Do(http.MethodPatch, "/api/v1/nodes/"+node.Name, strategicMerge, nodePatch, nil)
	if err != nil {
		return err
	}
	// get all pods on the node
	podList := &PodResourceList{}
	err = c.kubernetes.Do(http.MethodGet, "/api/v1/pods?fieldSelector=spec.nodeName%3D"+node.Name, nil, nil, podList)
	if err != nil {
		return err
	}
	log.Printf("Need to evict %d pods found on node '%s'\n", len(podList.Items), node.Name)
	// evict each pod
	errs := &catch.Errors{}
	for _, pod := range podList.Items {
		// TODO ignore pods owned by "apiVersion: apps/v1, kind: DaemonSet"
		log.Printf("Evicting pod '%s' in namespace '%s'\n", pod.Metadata.Name, pod.Metadata.Namespace)
		eviction := &PodEviction{
			ApiVersion: "policy/v1",
			Kind:       "Eviction",
			Metadata: PodEvictionMetadata{
				Name:      pod.Metadata.Name,
				Namespace: pod.Metadata.Namespace,
			},
		}
		errs.Add(c.kubernetes.Do(http.MethodPost, fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/eviction", pod.Metadata.Namespace, pod.Metadata.Name), rest.JSON(), eviction, nil))
	}
	return errs.OrNil()
}

func (c ClusterActions) GetNode(name string) (*NodeResource, error) {
	response := &NodeResource{}
	err := c.kubernetes.Do(http.MethodGet, "/api/v1/nodes/"+name, nil, nil, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c ClusterActions) DeleteNode(node NodeStatus) error {
	log.Printf("Deleting node %s with id %s\n", json.Format(node), node.CloudId)
	return c.kubernetes.Do(http.MethodDelete, "/api/v1/nodes/"+node.Name, nil, nil, nil)
}

func (c ClusterActions) GetClusterList() (*ClusterList, error) {
	var clusterList ClusterList
	err := c.kubernetes.Do(http.MethodGet, "/apis/hetzanetes.duncanpierce.org/v1/clusters", nil, nil, &clusterList)
	return &clusterList, err
}

func (c ClusterActions) SaveStatus(clusterName string, status *ClusterStatus) error {
	patch := Cluster{
		Status: status,
	}
	return c.kubernetes.Do(http.MethodPatch, "/apis/hetzanetes.duncanpierce.org/v1/clusters/"+clusterName+"/status", strategicMerge, patch, nil)
}

func (c ClusterActions) GetSshKeyIds() (keyIds []int, err error) {
	sshKeys := &hetzner.SshKeys{}
	err = c.hetzner.Do(http.MethodGet, "/ssh_keys", nil, nil, sshKeys)
	if err != nil {
		return
	}
	for _, sshKey := range sshKeys.SshKeys {
		keyIds = append(keyIds, sshKey.Id)
	}
	return
}
