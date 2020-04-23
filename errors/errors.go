// SPDX-License-Identifier: 0BSD

package errors

import "fmt"

// LsmError returned from JSON API
type LsmError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *LsmError) Error() string {
	if len(e.Data) > 0 {
		return fmt.Sprintf("code = %d, message = %s, data = %s", e.Code, e.Message, e.Data)
	}
	return fmt.Sprintf("code = %d, message = %s", e.Code, e.Message)
}

const (
	// LibBug ... Library bug
	LibBug = 1

	// PluginBug ... Bug found in plugin
	PluginBug = 2

	// JobStarted ... Job has been started
	JobStarted = 7

	// TimeOut ... Plugin timeout
	TimeOut = 11

	// DameonNotRunning ... lsmd does not appear to be running
	DameonNotRunning = 12

	// InvalidArgument ... provided argument is incorrect
	InvalidArgument = 101

	// PluginNotExist ... Plugin doesn't apprear to exist
	PluginNotExist = 311

	//TransPortComunication ... Issue reading/writing to plugin
	TransPortComunication = 400
)
