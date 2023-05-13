package kubeconfig

type (
	KubeConfig struct {
		ApiVersion     string    `json:"apiVersion"`
		Kind           string    `json:"kind"`
		CurrentContext string    `json:"current-context"`
		Contexts       []Context `json:"contexts"`
		Clusters       []Cluster `json:"clusters"`
		Users          []User    `json:"users"`
	}

	Context struct {
		Name    string `json:"name"`
		Context struct {
			Cluster string `json:"cluster"`
			User    string `json:"user"`
		} `json:"context"`
	}

	Cluster struct {
		Name    string `json:"name"`
		Cluster struct {
			CertificateAuthorityData string `json:"certificate-authority-data"`
			Server                   string `json:"server"`
		} `json:"cluster"`
	}

	User struct {
		Name string `json:"name"`
		User struct {
			Token string `json:"token"`
		} `json:"user"`
	}
)

func (k *KubeConfig) GetContext() *Context {
	for _, context := range k.Contexts {
		if context.Name == k.CurrentContext {
			return &context
		}
	}
	return nil
}

func (k *KubeConfig) GetClusterApiServer(name string) (server string, certificateBase64 string) {
	for _, cluster := range k.Clusters {
		if cluster.Name == name {
			return cluster.Cluster.Server, cluster.Cluster.CertificateAuthorityData
		}
	}
	return "", ""
}

func (k *KubeConfig) GetUserToken(name string) (token string) {
	for _, user := range k.Users {
		if user.Name == name {
			return user.User.Token
		}
	}
	return ""
}
