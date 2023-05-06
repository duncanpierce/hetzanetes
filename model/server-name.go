package model

import "fmt"

func GetServerName(clusterName string, nodeSetName string, generation int) string {
	return fmt.Sprintf("%s-%s-%d", clusterName, nodeSetName, generation)
}
