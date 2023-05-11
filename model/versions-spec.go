package model

func (v VersionsSpec) GetKubernetes() string {
	if v.Kubernetes == "" {
		return "stable"
	}
	return v.Kubernetes
}
