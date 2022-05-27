package util

import "strings"

// helper function to check existence in a collection, for some reason this doesnt exist in the go stdlib...
// if you want to use for other types just add to the generic parameter list
func Contains[T string | int | byte](collection []T, item T) bool {
	for _, value := range collection {
		if item == value {
			return true
		}
	}
	return false
}

func CleanFilePath(path string) string {
	return strings.ReplaceAll(strings.ReplaceAll(path, "/", "_"), "\\", "_")
}
