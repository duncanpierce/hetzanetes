package actions

import "github.com/duncanpierce/hetzanetes/model"

type (
	ClusterList struct {
		Items model.Clusters `json:"items"`
	}
)
