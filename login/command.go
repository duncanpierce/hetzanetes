package login

import (
	"fmt"
	"github.com/duncanpierce/hetzanetes/tmpl"
	"io/fs"
	"log"
	"path/filepath"
)

type (
	Command struct {
		Shell string
		Stdin []byte
	}
)

func CreateClusterCommands(clusterYaml []byte, clusterName, privateNetworkId, privateIpRange, k3sReleaseChannel, installHetzanetesVersion, hcloudToken, sshPrivateKey, sshPublicKey string) []Command {
	commands := []Command{
		{Shell: fmt.Sprintf("curl -sfL 'https://get.k3s.io' | INSTALL_K3S_CHANNEL=%s sh -s - server --cluster-init %s %s", k3sReleaseChannel, apiCommonConfig(), networkConfig(privateIpRange))},
		{Shell: fmt.Sprintf("kubectl create secret generic hcloud -n kube-system --from-literal=HCLOUD_TOKEN=%s --from-literal=HCLOUD_NETWORK=%s --from-literal=HCLOUD_NETWORK_ID=%s --from-literal=HCLOUD_NETWORK_IP_RANGE=%s", hcloudToken, clusterName, privateNetworkId, privateIpRange)},
		{Shell: fmt.Sprintf("kubectl create secret generic k3s -n kube-system --from-file=K3S_TOKEN=/var/lib/rancher/k3s/server/token")},
		{Shell: fmt.Sprintf("kubectl create secret generic ssh -n kube-system --from-literal=SSH_PRIVATE_KEY='%s' --from-literal=SSH_PUBLIC_KEY='%s'", sshPrivateKey, sshPublicKey)},
	}
	sendFileCommands, err := SendFiles(tmpl.Kustomize, "kustomize")
	if err != nil {
		panic(err)
	}
	commands = append(commands, sendFileCommands...)
	commands = append(commands,
		Command{Shell: "kubectl apply -k ."},
		Command{Stdin: []byte(tmpl.ClusterCrdYaml), Shell: "kubectl apply -f -"},
		Command{Stdin: []byte(tmpl.SetupRepairYaml), Shell: "kubectl apply -f -"},
	)
	if installHetzanetesVersion != "none" {
		repairYaml := fmt.Sprintf(tmpl.RepairClusterYaml, installHetzanetesVersion)
		commands = append(commands, Command{Stdin: []byte(repairYaml), Shell: "kubectl apply -f -"})
	}
	commands = append(commands, Command{Stdin: clusterYaml, Shell: "kubectl apply -f -"})
	return commands
}

func AddWorkerCommand(apiEndPoint, k3sJoinToken, k3sVersion, privateIpRange string) Command {
	shell := fmt.Sprintf("curl -sfL 'https://get.k3s.io' | K3S_URL=%s K3S_TOKEN=%s INSTALL_K3S_VERSION=%s sh -s - --kubelet-arg cloud-provider=external %s", apiEndPoint, k3sJoinToken, k3sVersion, networkConfig(privateIpRange))
	log.Printf("Add worker command: %s\n", shell)
	return Command{Shell: shell}
}

func AddApiServerCommand(apiEndPoint, k3sJoinToken, k3sVersion, privateIpRange string) Command {
	shell := fmt.Sprintf("curl -sfL 'https://get.k3s.io' | K3S_URL=%s K3S_TOKEN=%s INSTALL_K3S_VERSION=%s sh -s - server %s %s", apiEndPoint, k3sJoinToken, k3sVersion, apiCommonConfig(), networkConfig(privateIpRange))
	log.Printf("Add worker command: %s\n", shell)
	return Command{Shell: shell}
}

func SendFiles(filesystem fs.FS, directoryName string) ([]Command, error) {
	var commands []Command
	dirEntries, err := fs.ReadDir(filesystem, directoryName)
	if err != nil {
		panic(err)
	}
	for _, dirEntry := range dirEntries {
		fullPath := filepath.Join(directoryName, dirEntry.Name())
		bytes, err := fs.ReadFile(filesystem, fullPath)
		if err != nil {
			return nil, err
		}
		commands = append(commands, SendFile(dirEntry.Name(), bytes))
	}
	return commands, nil
}

func SendFile(filename string, contents []byte) Command {
	return Command{Stdin: contents, Shell: fmt.Sprintf("cat > '%s'", filename)}
}

func networkConfig(privateIpRange string) string {
	return fmt.Sprintf("--flannel-iface=$(ip -j route list %s | jq -r .[0].dev)", privateIpRange)
}

func apiCommonConfig() string {
	return "--disable servicelb --disable local-storage --disable-cloud-controller --kubelet-arg cloud-provider=external"
}
