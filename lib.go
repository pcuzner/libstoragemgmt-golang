package libstoragemgmt

import "os"

// UdsPath ... returns the unix domain file path
func udsPath() string {
	var p = os.Getenv(udsPathVarName)
	if len(p) > 0 {
		return p
	}
	return udsPathDefault
}
