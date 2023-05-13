package kubeconfig

import (
	"encoding/base64"
	"fmt"
	"github.com/ghodss/yaml"
	"os"
)

func InCluster() (server string, certificate []byte, token string, err error) {
	certificate, err = os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		return
	}
	var tokenBytes []byte
	tokenBytes, err = os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return
	}
	token = string(tokenBytes)
	server = "https://kubernetes.default.svc"
	return
}

func FromConfig(kubeconfig []byte) (server string, certificate []byte, token string, err error) {
	config := &KubeConfig{}
	err = yaml.Unmarshal(kubeconfig, &config)
	if err != nil {
		return "", nil, "", err
	}
	context := config.GetContext()
	if context == nil {
		return "", nil, "", fmt.Errorf("no context has been set in kubeconfig file")
	}
	var certificateBase64 string
	server, certificateBase64 = config.GetClusterApiServer(context.Context.Cluster)
	if server == "" || certificateBase64 == "" {
		return "", nil, "", fmt.Errorf("no server or certificate data has been set in kubeconfig file")
	}
	certificate, err = base64.StdEncoding.DecodeString(certificateBase64)
	if err != nil {
		return "", nil, "", err
	}
	token = config.GetUserToken(context.Context.User)
	if token == "" {
		return "", nil, "", fmt.Errorf("no auth token has been set in kubeconfig file")
	}
	return
}
