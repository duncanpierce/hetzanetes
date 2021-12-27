package k8s_client

import (
	"crypto/tls"
	"crypto/x509"
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
		Channel  string `json:"channel"`
		NodeSets `json:"nodeSets"`
	}
	NodeSets []NodeSet
	NodeSet  struct {
		Name      string   `json:"name"`
		ApiServer bool     `json:"apiServer"`
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

func (k *K8sClient) DoRaw(method string, path string) ([]byte, error) {
	request, err := http.NewRequest(method, k.BaseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+k.Token)
	response, err := k.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return body, err
}

func (k *K8sClient) Do(method string, path string, result interface{}) error {
	data, err := k.DoRaw(method, path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, result)
}

func (k *K8sClient) GetClusterList() (*ClusterList, error) {
	var clusterList ClusterList
	err := k.Do(http.MethodGet, "/apis/hetzanetes.duncanpierce.org/v1/clusters", &clusterList)
	return &clusterList, err
}
