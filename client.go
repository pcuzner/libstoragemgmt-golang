// SPDX-License-Identifier: 0BSD

package libstoragemgmt

import (
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

	var result string
	var libError = transport.invoke("plugin_register", args, &result)
	if libError != nil {
		return nil, libError
	}

	return &ClientConnection{tp: *transport, PluginName: pluginName, Timeout: timeout}, nil
}

// Close instructs the plugin to shutdown and exist.
func (c *ClientConnection) Close() error {
	var args = make(map[string]interface{})
	var result string
	var ourError = c.tp.invoke("plugin_unregister", args, &result)
	c.tp.close()
	return ourError
}

// Systems returns systems information
func (c *ClientConnection) Systems() ([]System, error) {
	var args = make(map[string]interface{})
	var systems []System
	var err = c.tp.invoke("systems", args, &systems)
	if err != nil {
		return systems, err
	}
	return systems, nil
}

// Volumes returns block device information
func (c *ClientConnection) Volumes() ([]Volume, error) {
	var args = make(map[string]interface{})
	var volumes []Volume
	var err = c.tp.invoke("volumes", args, &volumes)
	if err != nil {
		return volumes, err
	}
	return volumes, nil
}

// Pools returns the units of storage that block devices and FS
// can be created from.
func (c *ClientConnection) Pools() ([]Pool, error) {
	var args = make(map[string]interface{})
	var pools []Pool
	var err = c.tp.invoke("pools", args, &pools)
	if err != nil {
		return pools, err
	}
	return pools, nil
}

// Disks returns disks that are present.
func (c *ClientConnection) Disks() ([]Disk, error) {
	var args = make(map[string]interface{})
	var disks []Disk
	var err = c.tp.invoke("disks", args, &disks)
	if err != nil {
		return disks, err
	}
	return disks, nil
}

// FileSystems returns pools that are present.
func (c *ClientConnection) FileSystems() ([]FileSystem, error) {
	var args = make(map[string]interface{})
	var fileSystems []FileSystem
	var err = c.tp.invoke("fs", args, &fileSystems)
	if err != nil {
		return fileSystems, err
	}
	return fileSystems, nil
}

// NfsExports returns nfs exports  that are present.
func (c *ClientConnection) NfsExports() ([]NfsExport, error) {
	var args = make(map[string]interface{})
	var nfsExports []NfsExport
	var err = c.tp.invoke("exports", args, &nfsExports)
	if err != nil {
		return nfsExports, err
	}
	return nfsExports, nil
}

// AccessGroups returns access groups  that are present.
func (c *ClientConnection) AccessGroups() ([]AccessGroup, error) {
	var args = make(map[string]interface{})
	var accessGroups []AccessGroup
	var err = c.tp.invoke("access_groups", args, &accessGroups)
	if err != nil {
		return accessGroups, err
	}
	return accessGroups, nil
}

// TargetPorts returns target ports that are present.
func (c *ClientConnection) TargetPorts() ([]TargetPort, error) {
	var args = make(map[string]interface{})
	var targetPorts []TargetPort
	var err = c.tp.invoke("target_ports", args, &targetPorts)
	if err != nil {
		return targetPorts, err
	}
	return targetPorts, nil
}

// Batteries returns batteries that are present
func (c *ClientConnection) Batteries() ([]Battery, error) {
	var args = make(map[string]interface{})
	var batteries []Battery
	var err = c.tp.invoke("batteries", args, &batteries)
	if err != nil {
		return batteries, err
	}
	return batteries, nil
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
		var info []string
		var invokeError = trans.invoke("plugin_info", args, &info)

		trans.close()

		if invokeError != nil {
			return pluginInfos, invokeError
		}
		pluginInfos = append(pluginInfos, PluginInfo{
			Description: info[0],
			Version:     info[1],
			Name:        pluginPath})
	}

	return pluginInfos, nil
}
