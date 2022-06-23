package model

import "time"

type (
	CreateHetznerServerRequest struct {
		Name       string               `json:"name"`
		ServerType string               `json:"server_type"`
		Image      string               `json:"image"`
		Location   string               `json:"location"`
		Networks   []int                `json:"networks"`
		Firewalls  []HetznerFirewallRef `json:"firewalls"`
		Labels     map[string]string    `json:"labels"`
		SshKeys    []string             `json:"ssh_keys"`
		CloudInit  string               `json:"user_data"`
	}

	HetznerFirewallRef struct {
		Firewall int `json:"firewall"`
	}

	CreateHetznerServerResult struct {
		Server HetznerServerRef `json:"server"`
	}

	HetznerServerRef struct {
		Id         int              `json:"id"`
		PrivateNet []HetznerNetwork `json:"private_net"`
	}

	HetznerNetwork struct {
		Network int    `json:"network"`
		IP      string `json:"ip"`
	}

	HetznerServersResponse struct {
		Servers []HetznerServerResponse `json:"servers"`
	}

	HetznerServerResponse struct {
		Id         int       `json:"id"`
		Created    time.Time `json:"created"`
		Datacenter struct {
			Location struct {
				Name string `json:"name"`
			} `json:"location"`
		} `json:"datacenter"`
		Image struct {
			Name string `json:"name"`
		} `json:"image"`
		PrivateNets []HetznerNetwork `json:"private_net"`
		ServerType  struct {
			Name string `json:"name"`
		} `json:"server_type"`
	}
)
