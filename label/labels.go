package label

const (
	AppName             = "hetzanetes"
	ClusterNameLabel    = AppName + "-cluster"
	PrivateNetworkLabel = AppName + "-cluster-network"
	ApiServerLabel      = AppName + "-api"
	WorkerLabel         = AppName + "-worker"
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
