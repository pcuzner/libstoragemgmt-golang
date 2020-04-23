// SPDX-License-Identifier: 0BSD

package libstoragemgmt

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	errors "github.com/libstorage/libstoragemgmt-golang/errors"
)

// ClientConnection is the structure that encomposes the needed data for the plugin connection.
type ClientConnection struct {
	tp         transPort
	PluginName string
	Timeout    uint32
}

// Client establishes a connection to a plugin as specified in the URI.
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

	return &ClientConnection{tp: *transport, PluginName: pluginName, Timeout: timeout}, nil
}

// Close instructs the plugin to shutdown and exist.
func (c ClientConnection) Close() error {
	var args = make(map[string]interface{})
	var _, ourError = c.tp.invoke("plugin_unregister", args)
	c.tp.close()
	return ourError
}

// Systems returns systems information
func (c ClientConnection) Systems() ([]System, error) {
	var args = make(map[string]interface{})
	var systems []System
	var systemsJSON, err = c.tp.invoke("systems", args)
	if err != nil {
		return systems, err
	}
	var systemUnmarshal = json.Unmarshal(systemsJSON, &systems)
	if systemUnmarshal != nil {
		return systems, &errors.LsmError{
			Code:    errors.PluginBug,
			Message: fmt.Sprintf("Plugin returned unexpected system data %s", string(systemsJSON))}
	}

	return systems, nil
}

// AvailablePlugins retrieves all the available plugins.
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
