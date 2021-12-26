package k8s_client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
)

type (
	K8sClient struct {
		Client  http.Client
		BaseUrl string
		Token   string
	}
	ClusterList struct {
		Items ClusterItems `json:"items"`
	}
	ClusterItems []ClusterItem
	ClusterItem  struct {
		Metadata `json:"metadata"`
		Spec     `json:"spec"`
	}
	Metadata struct {
		Name string `json:"name"`
	}
	Spec struct {
		NodeSets `json:"nodeSets"`
	}
	NodeSets []NodeSet
	NodeSet  struct {
		Name      string   `json:"name"`
		Replicas  int      `json:"replicas"`
		NodeType  string   `json:"nodeType"`
		Locations []string `json:"locations"`
	}
)

func New() *K8sClient {
	cert, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
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
	tokenFile, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		panic(err)
	}
	return &K8sClient{
		BaseUrl: "https://kubernetes.default.svc",
		Client:  client,
		Token:   string(tokenFile),
	}
}

func (k *K8sClient) Do(method string, path string) (string, error) {
	request, err := http.NewRequest(method, k.BaseUrl+path, nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("Authorization", "Bearer "+k.Token)
	response, err := k.Client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return string(body), err
}

func (k *K8sClient) GetClusters() (string, error) {
	result, err := k.Do(http.MethodGet, "/apis/hetzanetes.duncanpierce.org/v1/clusters")
	fmt.Printf("Cluster JSON: %s\n", result)
	var clusterList ClusterList
	err = json.Unmarshal([]byte(result), &clusterList)
	if err != nil {
		return "", err
	}
	fmt.Printf("Struct: %#v\n", clusterList)
	return result, err
}
