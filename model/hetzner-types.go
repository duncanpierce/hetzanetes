package model

import "time"

type (
	CreateHetznerServerRequest struct {
		Name       string               `json:"name"`
		ServerType string               `json:"server_type"`
		Image      string               `json:"image"`
		Location   string               `json:"location,omitempty"`
		Networks   []int                `json:"networks,omitempty"`
		Firewalls  []HetznerFirewallRef `json:"firewalls,omitempty"`
		Labels     map[string]string    `json:"labels,omitempty"`
		SshKeys    []int                `json:"ssh_keys,omitempty"`
		CloudInit  string               `json:"user_data,omitempty"`
	}

	HetznerFirewallRef struct {
		Firewall int `json:"firewall"`
	}

	HetznerServerResult struct {
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

	HetznerSshKeys struct {
		SshKeys []HetznerSshKey `json:"ssh_keys,omitempty"`
	}

	HetznerSshKey struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		PublicKey string `json:"public_key,omitempty"`
	}

	HetznerActionsResponse struct {
		Actions []HetznerAction `json:"actions"`
	}

	HetznerAction struct {
		Command   string                  `json:"command"`
		Status    string                  `json:"status"`
		Error     HetznerActionError      `json:"error"`
		Started   time.Time               `json:"started"`
		Finished  time.Time               `json:"finished"`
		Resources []HetznerActionResource `json:"resources"`
	}

	HetznerActionError struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	HetznerActionResource struct {
		Id   int    `json:"id"`
		Type string `json:"type"`
	}
)
