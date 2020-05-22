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

func handleSearch(args map[string]interface{}, search []string) bool {
	var rc = true

	if len(search) == 0 {
		args["search_key"] = nil
		args["search_value"] = nil
	} else if len(search) == 2 {
		args["search_key"] = search[0]
		args["search_value"] = search[1]
	} else {
		rc = false
	}
	return rc
}
