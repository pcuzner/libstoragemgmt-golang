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

func handlePluginInfo(p *Plugin, params json.RawMessage) (interface{}, error) {
	return []string{p.desc, p.ver}, nil
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

func handleDisks(p *Plugin, params json.RawMessage) (interface{}, error) {
	return p.cb.San.Disks()
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

func handleVolumes(p *Plugin, params json.RawMessage) (interface{}, error) {
	var s search
	if uE := json.Unmarshal(params, &s); uE != nil {
		return nil, invalidArgs("volumes", uE)
	}

	if len(s.Key) > 0 {
		return p.cb.San.Volumes(s.Key, s.Value)
	}

	return p.cb.San.Volumes()
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

func exclusiveOr(item interface{}, job *string, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	}

	var result [2]interface{}

	if job != nil {
		result[0] = job
		result[1] = nil
	} else {
		result[0] = nil
		result[1] = item
	}
	return result, nil
}

func handleVolumeCreate(p *Plugin, params json.RawMessage) (interface{}, error) {

	type volumeCreateArgs struct {
		Pool         *Pool               `json:"pool"`
		Name         string              `json:"volume_name"`
		SizeBytes    uint64              `json:"size_bytes"`
		Provisioning VolumeProvisionType `json:"provisioning"`
		Flags        uint64              `json:"flags"`
	}

	var args volumeCreateArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_create", uE)
	}

	volume, jobID, error := p.cb.San.VolumeCreate(args.Pool, args.Name, args.SizeBytes, args.Provisioning)
	return exclusiveOr(volume, jobID, error)
}

func handleVolumeReplicate(p *Plugin, params json.RawMessage) (interface{}, error) {
	type volumeReplicateArgs struct {
		Pool    *Pool               `json:"pool"`
		RepType VolumeReplicateType `json:"rep_type"`
		Flags   uint64              `json:"flags"`
		SrcVol  Volume              `json:"volume_src"`
		Name    string              `json:"name"`
	}

	var args volumeReplicateArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_replicate", uE)
	}

	volume, jobID, error := p.cb.San.VolumeReplicate(args.Pool, args.RepType, &args.SrcVol, args.Name)
	return exclusiveOr(volume, jobID, error)
}

func handleVolumeReplicateRange(p *Plugin, params json.RawMessage) (interface{}, error) {
	type volumeReplicateRangeArgs struct {
		RepType VolumeReplicateType `json:"rep_type"`
		Ranges  []BlockRange        `json:"ranges"`
		SrcVol  Volume              `json:"volume_src"`
		DstVol  Volume              `json:"volume_dest"`
		Flags   uint64              `json:"flags"`
	}

	var a volumeReplicateRangeArgs
	if uE := json.Unmarshal(params, &a); uE != nil {
		return nil, invalidArgs("volume_replicate", uE)
	}

	return p.cb.San.VolumeReplicateRange(a.RepType, &a.SrcVol, &a.DstVol, a.Ranges)
}

func handleVolRepRangeBlockSize(p *Plugin, params json.RawMessage) (interface{}, error) {
	type args struct {
		System *System `json:"system"`
		Flags  uint64  `json:"flags"`
	}

	var a args
	if uE := json.Unmarshal(params, &a); uE != nil {
		return nil, invalidArgs("volume_replicate_range_block_size", uE)
	}
	return p.cb.San.VolumeRepRangeBlkSize(a.System)
}

func handleVolumeDelete(p *Plugin, params json.RawMessage) (interface{}, error) {
	type volumeDeleteArgs struct {
		Volume *Volume `json:"volume"`
		Flags  uint64  `json:"flags"`
	}

	var args volumeDeleteArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_delete", uE)
	}

	return p.cb.San.VolumeDelete(args.Volume)
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
		"plugin_info":                       handlePluginInfo,
		"plugin_register":                   nilAssign(c.Required.PluginRegister, handleRegister),
		"plugin_unregister":                 nilAssign(c.Required.PluginUnregister, handleUnRegister),
		"systems":                           nilAssign(c.Required.Systems, handleSystems),
		"capabilities":                      nilAssign(c.Required.Capabilities, handleCapabilities),
		"time_out_set":                      nilAssign(c.Required.TimeOutSet, handleTmoSet),
		"time_out_get":                      nilAssign(c.Required.TimeOutGet, handleTmoGet),
		"pools":                             nilAssign(c.Required.Pools, handlePools),
		"job_status":                        nilAssign(c.Required.JobStatus, handleJobStatus),
		"job_free":                          nilAssign(c.Required.JobFree, handleJobFree),
		"volume_create":                     nilAssign(c.San.VolumeCreate, handleVolumeCreate),
		"volume_delete":                     nilAssign(c.San.VolumeDelete, handleVolumeDelete),
		"volumes":                           nilAssign(c.San.Volumes, handleVolumes),
		"disks":                             nilAssign(c.San.Disks, handleDisks),
		"volume_replicate":                  nilAssign(c.San.VolumeReplicate, handleVolumeReplicate),
		"volume_replicate_range":            nilAssign(c.San.VolumeReplicateRange, handleVolumeReplicateRange),
		"volume_replicate_range_block_size": nilAssign(c.San.VolumeRepRangeBlkSize, handleVolRepRangeBlockSize),
	}
}
