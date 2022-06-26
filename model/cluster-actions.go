package model

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/rest"
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

var _ Actions = ClusterActions{}

func NewClusterClient() *ClusterActions {
	return &ClusterActions{
		kubernetes: NewKubernetes(),
		hetzner:    NewHetzner(env.HCloudToken()),
		k3s:        NewK3s(),
	}
}

func NewKubernetes() *rest.Client {
	cert, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		panic(err)
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
		panic(err)
	}
	return &rest.Client{
		BaseUrl: "https://kubernetes.default.svc",
		Http:    client,
		Token:   string(tokenFile),
	}
}

func NewHetzner(token string) *rest.Client {
	return &rest.Client{
		BaseUrl:     "https://api.hetzner.cloud/v1",
		Http:        http.Client{},
		Token:       token,
		LogResponse: true,
	}
}
func NewK3s() *rest.Client {
	return &rest.Client{
		BaseUrl: "https://update.k3s.io/v1-release",
		Http:    http.Client{},
	}
}

func (c ClusterActions) GetReleaseChannels() (ReleaseChannelStatuses, error) {
	response := &K3sReleaseChannelsResponse{}
	err := c.k3s.Do(http.MethodGet, "/channels", nil, nil, response)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c ClusterActions) GetBootstrapServer(name string, apiServer bool, kubernetesVersion *semver.Version) (*NodeStatus, error) {
	hetznerServers := &HetznerServersResponse{}
	err := c.hetzner.Do(http.MethodGet, "/servers?name="+name, nil, nil, hetznerServers)
	if err != nil {
		return nil, err
	}
	server := hetznerServers.Servers[0]
	network := server.PrivateNets[0]
	return &NodeStatus{
		Name:       name,
		ServerType: server.ServerType.Name,
		Location:   server.Datacenter.Location.Name,
		CloudId:    strconv.Itoa(server.Id),
		ClusterIP:  network.IP,
		BaseImage:  server.Image.Name,
		ApiServer:  apiServer,
		Version:    kubernetesVersion,
		Phases: PhaseChanges{
			PhaseChange{
				Phase:  Active,
				Reason: "bootstrap api server",
				Time:   server.Created,
			},
		},
	}, nil
}

func (c ClusterActions) CreateServer(name string, serverType string, image string, location string, privateNetworkId string, firewallIds []string, labels label.Labels, sshKeyIds []int, cloudInit string) (cloudId string, clusterIP string, err error) {
	privateNetworkNumber, _ := strconv.Atoi(privateNetworkId)
	firewalls := []HetznerFirewallRef{}
	for _, firewallId := range firewallIds {
		firewallNumber, _ := strconv.Atoi(firewallId)
		firewalls = append(firewalls, HetznerFirewallRef{firewallNumber})
	}

	serverRequest := CreateHetznerServerRequest{
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
	serverResult := &HetznerServerResult{}
	err = c.hetzner.Do(http.MethodPost, "/servers", rest.JSON(), serverRequest, serverResult)
	if err != nil {
		return
	}
	cloudIdNum := serverResult.Server.Id
	log.Printf("New server %s cloud id=%d\n", name, cloudIdNum)

	// Wait for private network to be attached (among other things)
	err = c.Await("servers", cloudIdNum)
	if err != nil {
		return
	}

	serverResult = &HetznerServerResult{}
	err = c.hetzner.Do(http.MethodGet, fmt.Sprintf("/servers/%d", cloudIdNum), nil, nil, serverResult)
	if err != nil {
		return
	}

	return strconv.Itoa(cloudIdNum), serverResult.Server.PrivateNet[0].IP, nil
}

func (c ClusterActions) Await(resourceType string, resourceId int) error {
	for {
		time.Sleep(1 * time.Second)
		result := &HetznerActionsResponse{}
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
		return c.hetzner.Do(http.MethodDelete, "/servers"+node.CloudId, nil, nil, nil) == rest.NotFound
	}
}

func (c ClusterActions) DrainNode(node NodeStatus) error {
	log.Printf("Draining node %#v\n", node)
	//TODO implement me
	return nil
}

func (c ClusterActions) GetKubernetesNode(node NodeStatus) (*NodeResource, error) {
	log.Printf("Checking node %#v ready\n", node)
	response := &NodeResource{}
	err := c.kubernetes.Do(http.MethodGet, "/api/v1/nodes/"+node.Name, nil, nil, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c ClusterActions) DeleteNode(node NodeStatus) error {
	log.Printf("Deleting node %#v with id %s\n", node, node.CloudId)
	return c.hetzner.Do(http.MethodDelete, "/servers/"+node.CloudId, nil, nil, nil)
}

func (c ClusterActions) GetClusterList() (*ClusterList, error) {
	var clusterList ClusterList
	err := c.kubernetes.Do(http.MethodGet, "/apis/hetzanetes.duncanpierce.org/v1/clusters", map[string]string{}, nil, &clusterList)
	return &clusterList, err
}

func (c ClusterActions) SaveStatus(clusterName string, status *ClusterStatus) error {
	patch := Cluster{
		Status: status,
	}
	headers := map[string]string{
		"Content-Type": "application/merge-patch+json",
	}
	return c.kubernetes.Do(http.MethodPatch, "/apis/hetzanetes.duncanpierce.org/v1/clusters/"+clusterName+"/status", headers, patch, nil)
}

func (c ClusterActions) GetSshKeyIds() (keyIds []int, err error) {
	sshKeys := &HetznerSshKeys{}
	err = c.hetzner.Do(http.MethodGet, "/ssh_keys", nil, nil, sshKeys)
	if err != nil {
		return
	}
	for _, sshKey := range sshKeys.SshKeys {
		keyIds = append(keyIds, sshKey.Id)
	}
	return
}
