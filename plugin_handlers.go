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

type search struct {
	Key   string `json:"search_key"`
	Value string `json:"search_value"`
	Flags uint64 `json:"flags"`
}

func handlePools(p *Plugin, params json.RawMessage) (interface{}, error) {
	var s search
	if uE := json.Unmarshal(params, &s); uE != nil {
		return nil, invalidArgs("pools", uE)
	}

	if len(s.Key) > 0 {
		return p.cb.Required.Pools(s.Key, s.Value)
	}

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

type jobArgs struct {
	ID string `json:"job_id"`
}

func handleJobStatus(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args jobArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("job_status", uE)
	}

	job, err := p.cb.Required.JobStatus(args.ID)
	if err != nil {
		return nil, err
	}

	var result [3]interface{}
	result[0] = job.Status
	result[1] = job.Percent
	result[2] = job.Item

	return result, nil
}

func handleJobFree(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args jobArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("job_status", uE)
	}

	return nil, p.cb.Required.JobFree(args.ID)
}

func handleVolumeCreate(p *Plugin, params json.RawMessage) (interface{}, error) {

	type volumeCreateArgs struct {
		Pool         *Pool               `json:"pool"`
		Name         string              `json:"volume_name"`
		SizeBytes    uint64              `json:"volume_size"`
		Provisioning VolumeProvisionType `json:"provisioning"`
		Flags        uint64              `json:"flags"`
	}

	var args volumeCreateArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_create", uE)
	}

	volume, jobID, error := p.cb.San.VolumeCreate(args.Pool, args.Name, args.SizeBytes, args.Provisioning)

	if error != nil {
		return nil, error
	}

	var result [2]interface{}
	result[0] = jobID
	result[1] = volume
	return result, nil
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
		"job_status":        nilAssign(c.Required.JobStatus, handleJobStatus),
		"job_free":          nilAssign(c.Required.JobFree, handleJobFree),
		"volume_create":     nilAssign(c.San.VolumeCreate, handleVolumeCreate),
	}
}
