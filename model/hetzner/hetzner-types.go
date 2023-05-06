package hetzner

import "time"

type (
	CreateServerRequest struct {
		Name       string            `json:"name"`
		ServerType string            `json:"server_type"`
		Image      string            `json:"image"`
		Location   string            `json:"location,omitempty"`
		Networks   []int             `json:"networks,omitempty"`
		Firewalls  []FirewallRef     `json:"firewalls,omitempty"`
		Labels     map[string]string `json:"labels,omitempty"`
		SshKeys    []int             `json:"ssh_keys,omitempty"`
		CloudInit  string            `json:"user_data,omitempty"`
	}

	FirewallRef struct {
		Firewall int `json:"firewall"`
	}

	ServerResult struct {
		Server ServerRef `json:"server"`
	}

	ServerRef struct {
		Id         int              `json:"id"`
		PrivateNet []PrivateNetwork `json:"private_net"`
	}

	PrivateNetwork struct {
		Network int    `json:"network"`
		IP      string `json:"ip"`
	}

	PublicNetwork struct {
		IPv4 PublicIp `json:"ipv4"`
		IPv6 PublicIp `json:"ipv6"`
	}

	PublicIp struct {
		IP string `json:"ip"`
	}

	ServersResponse struct {
		Servers []ServerResponse `json:"servers"`
	}

	ServerResponse struct {
		Id         int       `json:"id"`
		Name       string    `json:"name,omitempty"`
		Created    time.Time `json:"created"`
		Datacenter struct {
			Location struct {
				Name string `json:"name"`
			} `json:"location"`
		} `json:"datacenter"`
		Image struct {
			Name string `json:"name"`
		} `json:"image"`
		PublicNet   PublicNetwork    `json:"public_net,omitempty"`
		PrivateNets []PrivateNetwork `json:"private_net,omitempty"`
		ServerType  struct {
			Name string `json:"name"`
		} `json:"server_type"`
	}

	SshKeys struct {
		SshKeys []SshKey `json:"ssh_keys,omitempty"`
	}

	SshKey struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		PublicKey string `json:"public_key,omitempty"`
	}

	ActionsResponse struct {
		Actions []Action `json:"actions"`
	}

	Action struct {
		Command   string           `json:"command"`
		Status    string           `json:"status"`
		Error     ActionError      `json:"error"`
		Started   time.Time        `json:"started"`
		Finished  time.Time        `json:"finished"`
		Resources []ActionResource `json:"resources"`
	}

	ActionError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	ActionResource struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
	}
)
