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

func K3sToken() string {
	return Get("K3S_TOKEN", "a K3s join token")
}

func K3sEndpoint() string {
	return Get("K3S_URL", "a K3s API server endpoint")
}
