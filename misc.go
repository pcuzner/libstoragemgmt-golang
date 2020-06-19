// SPDX-License-Identifier: 0BSD

package libstoragemgmt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func getPluginIpcPath(pluginName string) string {
	return fmt.Sprintf("%s/%s", udsPath(), pluginName)
}

func getPlugins(path string) []string {
	var plugins []string

	// If we are walking the path when the daemon is starting we can get errors, loop
	// until we walk without errors or run out of time trying.  It is possible that
	// when we are walking the directory we only get a subset of the plugins that are
	// present.
	for i := 0; i < 10; i++ {

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
		} else {
			time.Sleep(time.Millisecond * 200)
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

// Callers should check for errors, then job id, then vol, but if they don't we want to make
// sure we we don't return something unexpected by mistake.
func ensureExclusiveVol(vol *Volume, job *string, err error) (*Volume, *string, error) {
	if err != nil {
		return nil, nil, err
	} else if job != nil {
		return nil, job, nil
	}
	return vol, nil, nil
}

// Callers should check for errors, then job id, then file system, but if they don't we want to make
// sure we we don't return something unexpected by mistake.
func ensureExclusiveFs(fs *FileSystem, job *string, err error) (*FileSystem, *string, error) {
	if err != nil {
		return nil, nil, err
	} else if job != nil {
		return nil, job, nil
	}
	return fs, nil, nil
}

// Callers should check for errors, then job id, then snapshot, but if they don't we want to make
// sure we we don't return something unexpected by mistake.
func ensureExclusiveSs(ss *FileSystemSnapShot, job *string, err error) (*FileSystemSnapShot, *string, error) {
	if err != nil {
		return nil, nil, err
	} else if job != nil {
		return nil, job, nil
	}
	return ss, nil, nil
}

// LsmBool is used to express booleans as we use 0 == false, 1 = true
// for the JSON RPC interface.
type LsmBool bool

// UnmarshalJSON used for custom JSON serialization
func (bit LsmBool) UnmarshalJSON(b []byte) error {
	bit = LsmBool(string(b) == "1")
	return nil
}

// MarshalJSON used to custom JSON serialization
func (bit LsmBool) MarshalJSON() ([]byte, error) {
	if bit {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

// MarshalJSON custom marshal for BlockRange
// ref. http://choly.ca/post/go-json-marshalling/
func (b *BlockRange) MarshalJSON() ([]byte, error) {
	type Alias BlockRange
	return json.Marshal(&struct {
		Class string `json:"class"`
		*Alias
	}{
		Class: "BlockRange",
		Alias: (*Alias)(b),
	})
}

