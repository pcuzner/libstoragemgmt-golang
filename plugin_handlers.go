// SPDX-License-Identifier: 0BSD

package libstoragemgmt

import (
	"encoding/json"
	"fmt"
	"reflect"

	errors "github.com/libstorage/libstoragemgmt-golang/errors"
)

func invalidArgs(msg string, e error) error {
	return &errors.LsmError{
		Code:    errors.TransPortInvalidArg,
		Message: fmt.Sprintf("%s: invalid arguments(s) %w\n", msg, e)}
}

func handleRegister(p *Plugin, params json.RawMessage) (interface{}, error) {

	var register PluginRegister
	if uE := json.Unmarshal(params, &register); uE != nil {
		return nil, invalidArgs("plugin_register", uE)
	}
	return nil, p.cb.Required.PluginRegister(&register)
}

func handleUnRegister(p *Plugin, params json.RawMessage) (interface{}, error) {
	return nil, p.cb.Required.PluginUnregister()
}

type tmoSetArgs struct {
	MS    uint32 `json:"ms"`
	Flags uint64 `json:"flags"`
}

func handleTmoSet(p *Plugin, params json.RawMessage) (interface{}, error) {
	var timeout tmoSetArgs
	if uE := json.Unmarshal(params, &timeout); uE != nil {
		return nil, invalidArgs("time_out_set", uE)
	}
	return nil, p.cb.Required.TimeOutSet(timeout.MS)
}

func handleTmoGet(p *Plugin, params json.RawMessage) (interface{}, error) {
	return p.cb.Required.TimeOutGet(), nil
}

func handleSystems(p *Plugin, params json.RawMessage) (interface{}, error) {
	return p.cb.Required.Systems()
}

func handlePools(p *Plugin, params json.RawMessage) (interface{}, error) {
	// TODO: Add search
	return p.cb.Required.Pools()
}

type capArgs struct {
	Sys System `json:"system"`
}

func handleCapabilities(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args capArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("time_out_set", uE)
	}
	return p.cb.Required.Capabilities(&args.Sys)
}

func nilAssign(present interface{}, cb handler) handler {

	// This seems like an epic fail of golang as I got burned by doing present == nil
	// ref. https://groups.google.com/forum/#!topic/golang-nuts/wnH302gBa4I/discussion
	if present == nil || reflect.ValueOf(present).IsNil() {
		return nil
	}
	return cb
}

func buildTable(c *CallBacks) map[string]handler {
	return map[string]handler{
		"plugin_register":   nilAssign(c.Required.PluginRegister, handleRegister),
		"plugin_unregister": nilAssign(c.Required.PluginUnregister, handleUnRegister),
		"systems":           nilAssign(c.Required.Systems, handleSystems),
		"capabilities":      nilAssign(c.Required.Capabilities, handleCapabilities),
		"time_out_set":      nilAssign(c.Required.TimeOutSet, handleTmoSet),
		"time_out_get":      nilAssign(c.Required.TimeOutGet, handleTmoGet),
		"pools":             nilAssign(c.Required.Pools, handlePools),
	}
}
