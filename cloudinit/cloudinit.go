package cloudinit

import (
	"bytes"
	"net"
	"text/template"
)

const cloudInitScript = `#cloud-config

# https://cloudinit.readthedocs.io/en/latest/topics/format.html

package_update: true
package_upgrade: true
packages:
  - ufw
  - curl
package_reboot_if_required: true
write_files:
  - path: /var/run/hetzanetes
    owner: root:root
    permissions: '0644'
    content: |
      # Hetzanetes test file

      Installed OK
runcmd:
  - ufw allow proto tcp from any to any port 22,6443
  - ufw allow from 10.244.0.0/16 # TODO try to remove this - should not be needed
  - ufw allow from 10.43.0.0/16
  - ufw allow from 10.42.0.0/16
  - ufw allow from 10.0.0.0/16
  - ufw -f default deny incoming
  - ufw -f default allow outgoing
  - ufw -f enable
  - "curl -sfL https://get.k3s.io | sh -s - --disable-cloud-controller --kubelet-arg cloud-provider=external"
  - "kubectl -n kube-system create secret generic hcloud --from-literal=token={{.ApiToken}} --from-literal=network={{.PrivateNetworkName}}"
  - "curl 'https://raw.githubusercontent.com/hetznercloud/hcloud-cloud-controller-manager/master/deploy/v1.7.0-networks.yaml' | awk '{sub(\"10\\.244\\.0\\.0/16\", \"10.42.0.0/16\"); print}' | kubectl apply -f -"
`

type ClusterConfig struct {
	ApiToken           string
	PrivateNetworkName string
	PrivateIpRange     net.IPNet
}

func CloudInitTemplate(config ClusterConfig) (string, error) {
	tmpl, err := template.New("cloudinit").Parse(cloudInitScript)
	if err != nil {
		return "", err
	}
	var yaml bytes.Buffer
	err = tmpl.Execute(&yaml, config)
	if err != nil {
		return "", err
	}
	return yaml.String(), nil
}
