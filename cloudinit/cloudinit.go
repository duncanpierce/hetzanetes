package cloudinit

import (
	"bytes"
	"fmt"
	"text/template"
)

const cloudInitScript = `#cloud-config

# https://cloudinit.readthedocs.io/en/latest/topics/format.html

package_update: true
package_upgrade: true
packages:
  - ufw
  - curl
  - jq
package_reboot_if_required: true
write_files:
  - path: /var/run/hetzanetes
    owner: root:root
    permissions: '0644'
    content: |
      # Hetzanetes test file

      Installed OK
runcmd:
  - ufw allow proto tcp from any to any port 22,6443 # TODO api servers only
  - ufw allow from 10.244.0.0/16 # TODO try to remove this - should not be needed
  - ufw allow from 10.43.0.0/16 # TODO try to parameterize
  - ufw allow from 10.42.0.0/16
  - ufw allow from 10.0.0.0/16
  - ufw -f default deny incoming
  - ufw -f default allow outgoing
  - ufw -f enable
  - "curl -sfL 'https://get.k3s.io' | {{ .K3sInstallEnvVars }} sh -s - {{ .K3sInstallArgs }} --kubelet-arg cloud-provider=external --flannel-iface=$(ip -j route list {{ .PrivateIpRange }} | jq -r .[0].dev)"
  - "kubectl -n kube-system create secret generic hcloud --from-literal=token={{.ApiToken}} --from-literal=network={{.PrivateNetworkName}}"
  - "curl -sfL 'https://raw.githubusercontent.com/hetznercloud/hcloud-cloud-controller-manager/master/deploy/v1.7.0-networks.yaml' | awk '{sub(\"10\\.244\\.0\\.0/16\", \"10.42.0.0/16\"); print}' | kubectl apply -f -"
`

// TODO there is also --flannel-conf : where is it and what's in it?

// TODO including the following causes collateral damage - should try to use something like https://github.com/krishicks/yaml-patch or Kustomize or send PR to make it config rather than hardcoded
//   - "curl 'https://raw.githubusercontent.com/hetznercloud/csi-driver/v1.4.0/deploy/kubernetes/hcloud-csi.yml' | awk '{sub(\"name: hcloud-csi\", \"name: hcloud\"); print}' | kubectl apply -f -"

type ClusterConfig struct {
	ApiToken           string
	PrivateNetworkName string
	PrivateIpRange     string
	K3sInstallEnvVars  string
	K3sInstallArgs     string
}

func Template(config ClusterConfig) (string, error) {
	tmpl, err := template.New("cloudinit").Parse(cloudInitScript)
	if err != nil {
		return "", err
	}
	var yaml bytes.Buffer
	err = tmpl.Execute(&yaml, config)
	if err != nil {
		return "", err
	}
	result := yaml.String()
	fmt.Print(result)
	return result, nil
}
