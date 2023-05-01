package model

func (v VersionsSpec) GetKubernetes() string {
	if v.Kubernetes == "" {
		return "stable"
	}
	return v.Kubernetes
}

func (v VersionsSpec) GetHetzanetes() string {
	if v.Hetzanetes == "" {
		return "latest"
	}
	return v.Hetzanetes
}
