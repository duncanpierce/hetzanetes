package model

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/duncanpierce/hetzanetes/env"
	"github.com/duncanpierce/hetzanetes/label"
	"github.com/duncanpierce/hetzanetes/rest"
	"net/http"
	"os"
	"strconv"
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

func (c ClusterActions) GetReleaseChannels() (r ReleaseChannelStatuses, err error) {
	err = c.k3s.Do(http.MethodGet, "/channels", nil, nil, &r)
	return
}

func (c ClusterActions) CreateServer(name string, serverType string, image string, location string, privateNetworkId string, firewallIds []string, labels label.Labels, sshKeys []string, cloudInit string) (cloudId string, err error) {
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
		SshKeys:    sshKeys,
		CloudInit:  cloudInit,
	}
	serverResult := CreateHetznerServerResult{}
	err = c.hetzner.Do(http.MethodPost, "/servers", rest.JSON(), serverRequest, serverResult)
	if err != nil {
		return
	}
	return strconv.Itoa(serverResult.Server.Id), nil
}

func (f ClusterActions) DeleteServer(cloudId string) (notFound bool) {
	return f.hetzner.Do(http.MethodDelete, "/servers"+cloudId, nil, nil, nil) == rest.NotFound
}

func (f ClusterActions) DrainNode(name string) error {
	//TODO implement me
	panic("implement me")
}

func (f ClusterActions) CheckNodeReady(name string) bool {
	//TODO implement me
	panic("implement me")
}

func (f ClusterActions) CheckNoNode(name string) bool {
	//TODO implement me
	panic("implement me")
}

func (f ClusterActions) DeleteNode(name string) {
	//TODO implement me
	panic("implement me")
}

func (f ClusterActions) GetClusterList() (*ClusterList, error) {
	var clusterList ClusterList
	err := f.kubernetes.Do(http.MethodGet, "/apis/hetzanetes.duncanpierce.org/v1/clusters", map[string]string{}, nil, &clusterList)
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
