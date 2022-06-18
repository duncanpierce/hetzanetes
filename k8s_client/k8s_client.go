package k8s_client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/duncanpierce/hetzanetes/model"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
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

func (k *K8sClient) DoRaw(method string, path string, headers map[string]string, requestBody []byte) ([]byte, error) {
	request, err := http.NewRequest(method, k.BaseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+k.Token)
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	request.Body = io.NopCloser(bytes.NewReader(requestBody))
	response, err := k.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if response.StatusCode >= 400 {
		return responseBody, fmt.Errorf("got status code %d from Kubernetes API", response.StatusCode)
	}
	return responseBody, err
}

func (k *K8sClient) Do(method string, path string, headers map[string]string, request interface{}, result interface{}) error {
	var requestBody []byte
	var err error
	var data []byte

	if request != nil {
		requestBody, err = json.Marshal(request)
		log.Printf("formatted request body as %s\n", string(requestBody))
	}
	if err != nil {
		return err
	}
	data, err = k.DoRaw(method, path, headers, requestBody)
	if err != nil {
		return err
	}
	if result != nil {
		return json.Unmarshal(data, result)
	}
	return nil
}

func (k *K8sClient) GetClusterList() (*ClusterList, error) {
	var clusterList ClusterList
	err := k.Do(http.MethodGet, "/apis/hetzanetes.duncanpierce.org/v1/clusters", map[string]string{}, nil, &clusterList)
	return &clusterList, err
}

func (k *K8sClient) UpdateStatus(clusterName string, status *model.Status) error {
	patch := model.Cluster{
		Status: *status,
	}
	headers := map[string]string{
		"Content-Type": "application/merge-patch+json",
	}
	err := k.Do(http.MethodPatch, "/apis/hetzanetes.duncanpierce.org/v1/clusters/"+clusterName+"/status", headers, patch, nil)
	return err
}
