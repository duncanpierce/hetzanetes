package actions

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
		Id         int                 `json:"id"`
		PrivateNet []HetznerNetworkRef `json:"private_net"`
	}

	HetznerNetworkRef struct {
		Network int    `json:"network"`
		IP      string `json:"ip"`
	}
)
