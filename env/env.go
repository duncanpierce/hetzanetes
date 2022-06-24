package env

import (
	"fmt"
	"os"
)

func Get(name string, description string) string {
	value := os.Getenv(name)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s must contain %s", name, description))
	}
	return value
}

func HCloudToken() string {
	return Get("HCLOUD_TOKEN", "a Hetzner Cloud API token")
}

func HCloudNetwork() string {
	return Get("HCLOUD_NETWORK", "a Hetzner cloud network name")
}

func HCloudNetworkId() string {
	return Get("HCLOUD_NETWORK_ID", "a Hetzner cloud network id")
}

func HCloudNetworkIpRange() string {
	return Get("HCLOUD_NETWORK_IP_RANGE", "IP range of cluster network")
}

func K3sToken() string {
	return Get("K3S_TOKEN", "a K3s join token")
}

func K3sEndpoint() string {
	return Get("K3S_URL", "a K3s API server endpoint")
}

func SshPrivateKey() string {
	return Get("SSH_PRIVATE_KEY", "the SSH private key for the cluster")
}

func SshPublicKey() string {
	return Get("SSH_PUBLIC_KEY", "the SSH public key for the cluster")
}
