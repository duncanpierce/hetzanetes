package tmpl

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

//go:embed files/*
var files embed.FS

type ClusterConfig struct {
	ApiEndpoint        string
	HetznerApiToken    string
	PrivateNetworkName string
	PrivateIpRange     string // TODO define IpRange map[string]string and read as {{.IpRange.PrivateNetwork}} etc - maybe rename to ClusterNetwork
	PodIpRange         string
	ServiceIpRange     string
	InstallDirectory   string
	JoinToken          string
	ServerType         string
	// TODO add Version map[string]string and emit versions in files unless the key is missing
}

func Template(config ClusterConfig, templateName string) string {
	t, err := template.ParseFS(files, "files/*")
	if err != nil {
		panic(fmt.Sprintf("error loading templates: %s", err.Error()))
	}
	var buffer bytes.Buffer
	err = t.ExecuteTemplate(&buffer, templateName, config)
	if err != nil {
		panic(fmt.Sprintf("error expanding template: %s", err.Error()))
	}
	result := buffer.String()
	return result
}
