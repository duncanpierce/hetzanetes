package tmpl

import (
	"embed"
)

//go:embed kustomize/*
var Kustomize embed.FS

//go:embed default-cluster.yaml
var DefaultClusterYaml string

//go:embed cluster-crd.yaml
var ClusterCrdYaml string

//go:embed repair-cluster.yaml
var RepairClusterYaml string
