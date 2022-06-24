package tmpl

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/Masterminds/sprig"
	"os"
	"text/template"
)

//go:embed cloudinit/*
var cloudinit embed.FS

//go:embed kustomize/*
var kustomize embed.FS

//go:embed default-cluster.yaml
var defaultCluster string

type ClusterConfig struct {
	ApiEndpoint       string
	HetznerApiToken   string
	ClusterName       string
	PrivateIpRange    string // TODO define IpRange map[string]string and read as {{.IpRange.PrivateNetwork}} etc - maybe rename to ClusterNetwork
	ClusterNetworkId  string
	PodIpRange        string
	ServiceIpRange    string
	InstallDirectory  string
	JoinToken         string
	K3sReleaseChannel string
	KubernetesVersion string
	HetzanetesTag     string
	ClusterYaml       string
	SshPublicKey      string
	SshPrivateKey     string
	// TODO add Version map[string]string and emit versions in files unless the key is missing
}

func Cloudinit(config ClusterConfig, templateName string) string {
	return Expand(Parse(cloudinit, "cloudinit"), templateName, config)
}

func Parse(files embed.FS, directory string) *template.Template {
	t, err := template.New("template").Funcs(sprig.TxtFuncMap()).ParseFS(files, directory+"/*")
	if err != nil {
		panic(fmt.Sprintf("error loading templates: %s", err.Error()))
	}
	return t
}

func Expand(t *template.Template, templateName string, config ClusterConfig) string {
	var buffer bytes.Buffer
	err := t.ExecuteTemplate(&buffer, templateName, config)
	if err != nil {
		panic(fmt.Sprintf("error expanding template: %s", err.Error()))
	}
	result := buffer.String()
	return result
}

func WriteKustomizeFiles(config ClusterConfig) error {
	for _, template := range Parse(kustomize, "kustomize").Templates() {
		file, err := os.Create(template.Name())
		if err != nil {
			return err
		}
		defer file.Close()
		template.Execute(file, config)
	}
	return nil
}

func DefaultClusterFile(clusterName string) ([]byte, error) {
	t, err := template.New("default").Parse(defaultCluster)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	err = t.Execute(b, ClusterConfig{ClusterName: clusterName})
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
