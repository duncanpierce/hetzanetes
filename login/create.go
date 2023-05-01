package login

import "fmt"

var (
	apiCommonConfig     = "--disable servicelb --disable local-storage --disable-cloud-controller --kubelet-arg cloud-provider=external"
	getPrivateInterface = "$(ip -j route list {{.PrivateIpRange}} | jq -r .[0].dev)"
)

func CreateCommands(clusterName, privateNetworkId, privateIpRange, k3sReleaseChannel, hcloudToken, sshPrivateKey, sshPublicKey string) (commands []string, env map[string]string) {
	return []string{
			fmt.Sprintf("curl -sfL 'https://get.k3s.io' | INSTALL_K3S_CHANNEL=%s sh -s - server --cluster-init %s --flannel-iface=%s", k3sReleaseChannel, apiCommonConfig, getPrivateInterface),
			fmt.Sprintf("kubectl create secret generic hcloud -n kube-system --from-literal=HCLOUD_TOKEN=%s --from-literal=HCLOUD_NETWORK=%s --from-literal=HCLOUD_NETWORK_ID=%s --from-literal=HCLOUD_NETWORK_IP_RANGE=%s", hcloudToken, clusterName, privateNetworkId, privateIpRange),
			"kubectl create secret generic k3s -n kube-system --from-file=K3S_TOKEN=/var/lib/rancher/k3s/server/token",
			"kubectl create secret generic ssh -n kube-system --from-literal=SSH_PRIVATE_KEY=\"$LC_SSH_PRIVATE_KEY\" --from-literal=SSH_PUBLIC_KEY=\"$LC_SSH_PUBLIC_KEY\"",
		}, map[string]string{
			// Many SSHD servers have an allow-list of env var names configured, LC_* is commonly allowed
			"LC_SSH_PRIVATE_KEY": sshPrivateKey,
			"LC_SSH_PUBLIC_KEY":  sshPublicKey,
		}
}

func AddWorkerCommands(apiEndPoint, k3sJoinToken, k3sVersion string) (commands []string) {
	return []string{
		fmt.Sprintf("curl -sfL 'https://get.k3s.io' | K3S_URL=%s K3S_TOKEN=%s INSTALL_K3S_VERSION=%s sh -s - --kubelet-arg cloud-provider=external --flannel-iface=%s", apiEndPoint, k3sJoinToken, k3sVersion, getPrivateInterface),
	}
}

func AddApiServerCommands(apiEndPoint, k3sJoinToken, k3sVersion string) (commands []string) {
	return []string{
		fmt.Sprintf("curl -sfL 'https://get.k3s.io' | K3S_URL=%s K3S_TOKEN=%s INSTALL_K3S_VERSION=%s sh -s - server %s --flannel-iface=%s", apiEndPoint, k3sJoinToken, k3sVersion, apiCommonConfig, getPrivateInterface),
	}
}
