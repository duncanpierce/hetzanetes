package tmpl

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed cloudinit/*
var cloudinit embed.FS

//go:embed kustomize/*
var kustomize embed.FS

type ClusterConfig struct {
	ApiEndpoint      string
	HetznerApiToken  string
	ClusterName      string
	PrivateIpRange   string // TODO define IpRange map[string]string and read as {{.IpRange.PrivateNetwork}} etc - maybe rename to ClusterNetwork
	PodIpRange       string
	ServiceIpRange   string
	InstallDirectory string
	JoinToken        string
	ServerType       string
	// TODO add Version map[string]string and emit versions in files unless the key is missing
}

func Cloudinit(config ClusterConfig, templateName string) string {
	return Expand(Parse(cloudinit, "cloudinit"), templateName, config)
}

func Parse(files embed.FS, directory string) *template.Template {
	template, err := template.ParseFS(files, directory+"/*")
	if err != nil {
		panic(fmt.Sprintf("error loading templates: %s", err.Error()))
	}
	return template
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
