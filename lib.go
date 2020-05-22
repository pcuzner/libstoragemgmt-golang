package libstoragemgmt

import (
	"os"
)

// UdsPath ... returns the unix domain file path
func udsPath() string {
	var p = os.Getenv(udsPathVarName)
	if len(p) > 0 {
		return p
	}
	return udsPathDefault
}

func contains(s []string, v string) bool {
	for _, a := range s {
		if a == v {
			return true
		}
	}
	return false
}

func emptySliceIfNil(provided []string) []string {
	var empty = make([]string, 0)

	if provided != nil {
		return provided
	}
	return empty
}
