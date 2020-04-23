// SPDX-License-Identifier: 0BSD

package libstoragemgmt

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	errors "github.com/libstorage/libstoragemgmt-golang/errors"
)

// ClientConnection ... structure for client connection
type ClientConnection struct {
	tp         transPort
	pluginName string
	timeout    uint32
}

// Client ... Establish a connection with the specified plugin URI and password
func Client(uri string, password string, timeout uint32) (*ClientConnection, error) {

	var p, parseError = url.Parse(uri)
	if parseError != nil {
		return nil, &errors.LsmError{
			Code:    errors.InvalidArgument,
			Message: fmt.Sprintf("invalid uri: %w", parseError)}
	}

	var pluginName = p.Scheme
	var pluginIpcPath = getPluginIpcPath(pluginName)

	var transport, transPortError = newTransport(pluginIpcPath, true)
	if transPortError != nil {
		return nil, transPortError
	}

	var args = make(map[string]interface{})
	args["password"] = password
	args["uri"] = uri
	args["timeout"] = timeout

	var _, libError = transport.invoke("plugin_register", args)
	if libError != nil {
		return nil, libError
	}

	return &ClientConnection{tp: *transport, pluginName: pluginName, timeout: timeout}, nil
}

// Close the plugin
func (c ClientConnection) Close() error {
	var args = make(map[string]interface{})
	var _, ourError = c.tp.invoke("plugin_unregister", args)
	c.tp.close()
	return ourError
}

func getPluginIpcPath(pluginName string) string {
	return fmt.Sprintf("%s/%s", udsPath(), pluginName)
}

// PluginInfo - Information about a specific plugin
type PluginInfo struct {
	Version     string
	Description string
	Name        string
}

func getPlugins(path string) []string {
	var plugins []string

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

// AvailablePlugins retrieves all the available plugins
func AvailablePlugins() ([]PluginInfo, error) {
	var udsPath = udsPath()

	if _, err := os.Stat(udsPath); os.IsNotExist(err) {
		return make([]PluginInfo, 0), &errors.LsmError{
			Code:    errors.DameonNotRunning,
			Message: fmt.Sprintf("LibStorageMgmt daemon is not running for socket folder %s", udsPath),
			Data:    ""}
	}

	var pluginInfos []PluginInfo
	for _, pluginPath := range getPlugins(udsPath) {

		var trans, transError = newTransport(pluginPath, true)
		if transError != nil {
			return nil, transError
		}

		var args = make(map[string]interface{})
		var reply, invokeError = trans.invoke("plugin_info", args)

		trans.close()

		if invokeError != nil {
			return pluginInfos, invokeError
		}

		var info []string
		var infoError = json.Unmarshal(reply, &info)
		if infoError != nil {

		}
		pluginInfos = append(pluginInfos, PluginInfo{Description: info[0], Version: info[1], Name: pluginPath})
	}

	return pluginInfos, nil
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
