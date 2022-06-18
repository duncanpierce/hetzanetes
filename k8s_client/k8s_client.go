package k8s_client

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/duncanpierce/hetzanetes/model"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"os"
)

type (
	K8sClient struct {
		Client  http.Client
		BaseUrl string
		Token   string
	}
	ClusterList struct {
		Items model.Clusters `json:"items"`
	}
)

func New() *K8sClient {
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
	body, err := io.ReadAll(response.Body)
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
