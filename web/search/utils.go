package search

import (
	"diskhub/web/models"
	"slices"
)

type KeyValue struct {
	Project *models.Project
	Value   int
}

func RemoveString(list []string, s string) []string {
	// Get the pos of the string
	for i, obj := range list {
		if obj == s {
			// Remove the string
			return slices.Delete(list, i, i+1)
		}
	}
	return list
}
