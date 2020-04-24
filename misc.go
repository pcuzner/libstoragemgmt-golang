// SPDX-License-Identifier: 0BSD

package libstoragemgmt

import (
	"fmt"
	"os"
	"path/filepath"
)

func getPluginIpcPath(pluginName string) string {
	return fmt.Sprintf("%s/%s", udsPath(), pluginName)
}

func getPlugins(path string) []string {
	var plugins []string

	// TODO: Put a limit on trying here.  Idea is we could get errors
	// while the daemon is starting, but we shouldn't loop forever.
	for true {

		plugins = nil

		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() == false {
				plugins = append(plugins, path)
			}
			return nil
		})

		if err == nil {
			break
		}
	}
	return plugins
}

func checkDaemonExists() bool {
	var present = false
	var udsPath = udsPath()

	// The unix domain socket needs to exist
	if _, err := os.Stat(udsPath); os.IsNotExist(err) {
		return present
	}

	for _, pluginPath := range getPlugins(udsPath) {
		var trans, err = newTransport(pluginPath, false)
		if err == nil {
			present = true
			trans.close()
		}
		break
	}

	return present
}

// LsmBool is used to express booleans as we use 0 == false, 1 = true
// for the JSON RPC interface.
type LsmBool bool

// UnmarshalJSON used for custom JSON serialization
func (bit *LsmBool) UnmarshalJSON(b []byte) error {
	*bit = LsmBool(string(b) == "1")
	return nil
}
