package tmpl

import (
	"embed"
)

//go:embed kustomize/*
var Kustomize embed.FS

//go:embed default-cluster.yaml
var DefaultCluster string
