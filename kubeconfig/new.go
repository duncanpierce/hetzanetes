package kubeconfig

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func InCluster() (server string, certificate []byte, token []byte, err error) {
	certificate, err = os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		return
	}
	token, err = os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return
	}
	server = "https://kubernetes.default.svc"
	return
}

func FromConfig(kubeconfig []byte) (server string, certificate []byte, token []byte, err error) {
	config := &KubeConfig{}
	err = yaml.Unmarshal(kubeconfig, &config)
	if err != nil {
		return "", nil, nil, err
	}
	context := config.GetContext()
	if context == nil {
		return "", nil, nil, fmt.Errorf("no context has been set in kubeconfig file")
	}
	server, certificate = config.GetClusterApiServer(context.Context.User)
	if server == "" || certificate == nil {
		return "", nil, nil, fmt.Errorf("no server or certificate data has been set in kubeconfig file")
	}
	token = config.GetUserToken(context.Context.User)
	if token == nil {
		return "", nil, nil, fmt.Errorf("no auth token has been set in kubeconfig file")
	}
	return
}
