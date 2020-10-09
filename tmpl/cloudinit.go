package tmpl

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type ClusterConfig struct {
	ApiServer          bool
	ApiEndpoint        string
	HetznerApiToken    string
	PrivateNetworkName string
	PrivateIpRange     string // TODO define IpRange map[string]string and read as {{.IpRange.PrivateNetwork}} etc - maybe rename to ClusterNetwork
	PodIpRange         string
	ServiceIpRange     string
	InstallDirectory   string
	// TODO add Version map[string]string and emit versions in templates unless the key is missing
}

func Template(config ClusterConfig) string {
	template := template.New("")
	for _, assetName := range AssetNames() {
		assetNameAndExtension := strings.Split(assetName, ".")
		tmplSource := "{{define \"" + assetNameAndExtension[0] + "\"}}" + string(MustAsset(assetName)) + "{{end}}"
		_, err := template.Parse(tmplSource)
		if err != nil {
			panic(fmt.Sprintf("error parsing template %s", assetName))
		}
	}
	var buffer bytes.Buffer
	err := template.ExecuteTemplate(&buffer, "cloudinit", config)
	if err != nil {
		panic(fmt.Sprintf("error expanding template: %s", err.Error()))
	}
	return buffer.String()
}
