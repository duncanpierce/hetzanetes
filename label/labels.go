package label

const (
	AppName        = "hetzanetes"
	Cluster        = AppName + "-cluster"
	PrivateNetwork = AppName + "-cluster-network"
	ApiServer      = AppName + "-api"
	Worker         = AppName + "-worker"
	NodeSet        = AppName + "-nodeset"
)

type Labels map[string]string

func (original Labels) Copy() Labels {
	copy := Labels{}
	for k, v := range original {
		copy[k] = v
	}
	return copy
}

func (l Labels) Mark(key string) Labels {
	return l.Set(key, "")
}

func (l Labels) Set(key, value string) Labels {
	l[key] = value
	return l
}
