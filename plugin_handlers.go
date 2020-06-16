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
	return nil, p.cb.Mgmt.PluginRegister(&register)
}

func handleUnRegister(p *Plugin, params json.RawMessage) (interface{}, error) {
	return nil, p.cb.Mgmt.PluginUnregister()
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
	return nil, p.cb.Mgmt.TimeOutSet(timeout.MS)
}

func handleTmoGet(p *Plugin, params json.RawMessage) (interface{}, error) {
	return p.cb.Mgmt.TimeOutGet(), nil
}

func handleSystems(p *Plugin, params json.RawMessage) (interface{}, error) {
	return p.cb.Mgmt.Systems()
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
		return p.cb.Mgmt.Pools(s.Key, s.Value)
	}

	return p.cb.Mgmt.Pools()
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
	return p.cb.Mgmt.Capabilities(&args.Sys)
}

type jobArgs struct {
	ID string `json:"job_id"`
}

func handleJobStatus(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args jobArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("job_status", uE)
	}

	job, err := p.cb.Mgmt.JobStatus(args.ID)
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

	return nil, p.cb.Mgmt.JobFree(args.ID)
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

func handleVolumeResize(p *Plugin, params json.RawMessage) (interface{}, error) {
	type args struct {
		Volume *Volume `json:"volume"`
		Size   uint64  `json:"new_size_bytes"`
		Flags  uint64  `json:"flags"`
	}

	var a args
	if uE := json.Unmarshal(params, &a); uE != nil {
		return nil, invalidArgs("volume_resize", uE)
	}

	fmt.Printf("args = %+v\n", a)

	volume, jobID, error := p.cb.San.VolumeResize(a.Volume, a.Size)
	return exclusiveOr(volume, jobID, error)
}

type volumeArgument struct {
	Volume *Volume `json:"volume"`
	Flags  uint64  `json:"flags"`
}

func handleVolumeEnable(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args volumeArgument
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_enable", uE)
	}

	return nil, p.cb.San.VolumeEnable(args.Volume)
}

func handleVolumeDisable(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args volumeArgument
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_disable", uE)
	}

	return nil, p.cb.San.VolumeDisable(args.Volume)
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

type maskArgs struct {
	Vol   *Volume      `json:"volume"`
	Ag    *AccessGroup `json:"access_group"`
	Flags uint64       `json:"flags"`
}

func handleVolumeMask(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args maskArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_mask", uE)
	}

	return nil, p.cb.San.VolumeMask(args.Vol, args.Ag)
}

func handleVolumeUnMask(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args maskArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_unmask", uE)
	}

	return nil, p.cb.San.VolumeUnMask(args.Vol, args.Ag)
}

func handleVolsMaskedToAg(p *Plugin, params json.RawMessage) (interface{}, error) {
	type argsAg struct {
		Ag    *AccessGroup `json:"access_group"`
		Flags uint64       `json:"flags"`
	}

	var args argsAg
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volumes_accessible_by_access_group", uE)
	}

	return p.cb.San.VolsMaskedToAg(args.Ag)
}

func handleAccessGroups(p *Plugin, params json.RawMessage) (interface{}, error) {
	return p.cb.San.AccessGroups()
}

func handleAccessGroupCreate(p *Plugin, params json.RawMessage) (interface{}, error) {
	type agCreateArgs struct {
		Name     string        `json:"name"`
		InitID   string        `json:"init_id"`
		InitType InitiatorType `json:"init_type"`
		System   *System       `json:"system"`
		Flags    uint64        `json:"flags"`
	}

	var args agCreateArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("access_group_create", uE)
	}

	return p.cb.San.AccessGroupCreate(args.Name, args.InitID, args.InitType, args.System)
}

func handleAccessGroupDelete(p *Plugin, params json.RawMessage) (interface{}, error) {
	type agDeleteArgs struct {
		Ag    *AccessGroup `json:"access_group"`
		Flags uint64       `json:"flags"`
	}

	var args agDeleteArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("access_group_delete", uE)
	}

	return nil, p.cb.San.AccessGroupDelete(args.Ag)
}

type accessGroupInitArgs struct {
	Ag       *AccessGroup  `json:"access_group"`
	ID       string        `json:"init_id"`
	InitType InitiatorType `json:"init_type"`
	Flags    uint64        `json:"flags"`
}

func handleAccessGroupInitAdd(p *Plugin, params json.RawMessage) (interface{}, error) {

	var args accessGroupInitArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("access_group_initiator_add", uE)
	}

	return p.cb.San.AccessGroupInitAdd(args.Ag, args.ID, args.InitType)
}

func handleAccessGroupInitDelete(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args accessGroupInitArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("access_group_initiator_delete", uE)
	}

	return p.cb.San.AccessGroupInitDelete(args.Ag, args.ID, args.InitType)
}

func handleAgsGrantedToVol(p *Plugin, params json.RawMessage) (interface{}, error) {
	type argsVol struct {
		Vol   *Volume `json:"volume"`
		Flags uint64  `json:"flags"`
	}

	var args argsVol
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("access_groups_granted_to_volume", uE)
	}

	return p.cb.San.AgsGrantedToVol(args.Vol)
}

func handleIscsiChapAuthSet(p *Plugin, params json.RawMessage) (interface{}, error) {
	type argsIscsi struct {
		InitID      string  `json:"init_id"`
		InUser      *string `json:"in_user"`
		InPassword  *string `json:"in_password"`
		OutUser     *string `json:"out_user"`
		OutPassword *string `json:"out_password"`
		Flags       uint64  `json:"flags"`
	}

	var args argsIscsi
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("iscsi_chap_auth", uE)
	}

	return nil, p.cb.San.IscsiChapAuthSet(args.InitID, args.InUser, args.InPassword, args.OutUser, args.OutPassword)
}

type argsChildDep struct {
	Vol   *Volume `json:"volume"`
	Flags uint64  `json:"flags"`
}

func handleVolHasChildDep(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args argsChildDep
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_child_dependency", uE)
	}

	return p.cb.San.VolHasChildDep(args.Vol)
}

func handleVolChildDepRm(p *Plugin, params json.RawMessage) (interface{}, error) {
	var args argsChildDep
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("volume_child_dependency_rm", uE)
	}

	return p.cb.San.VolChildDepRm(args.Vol)
}

func handleTargetPorts(p *Plugin, params json.RawMessage) (interface{}, error) {
	return p.cb.San.TargetPorts()
}

func handleFs(p *Plugin, params json.RawMessage) (interface{}, error) {
	var s search
	if uE := json.Unmarshal(params, &s); uE != nil {
		return nil, invalidArgs("fs", uE)
	}

	if len(s.Key) > 0 {
		return p.cb.File.FileSystems(s.Key, s.Value)
	}
	return p.cb.File.FileSystems()
}

func handleFsCreate(p *Plugin, params json.RawMessage) (interface{}, error) {
	type fsCreateArgs struct {
		Pool      *Pool  `json:"pool"`
		Name      string `json:"name"`
		SizeBytes uint64 `json:"size_bytes"`
		Flags     uint64 `json:"flags"`
	}

	var args fsCreateArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_create", uE)
	}

	fs, jobID, error := p.cb.File.FsCreate(args.Pool, args.Name, args.SizeBytes)
	return exclusiveOr(fs, jobID, error)
}

func handleFsDelete(p *Plugin, params json.RawMessage) (interface{}, error) {
	type fsDeleteArgs struct {
		Fs    *FileSystem `json:"fs"`
		Flags uint64      `json:"flags"`
	}

	var args fsDeleteArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_delete", uE)
	}
	return p.cb.File.FsDelete(args.Fs)
}

func handleFsResize(p *Plugin, params json.RawMessage) (interface{}, error) {
	type fsResizeArgs struct {
		Fs    *FileSystem `json:"fs"`
		Size  uint64      `json:"new_size_bytes"`
		Flags uint64      `json:"flags"`
	}

	var args fsResizeArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_resize", uE)
	}

	fs, job, err := p.cb.File.FsResize(args.Fs, args.Size)
	return exclusiveOr(fs, job, err)

}

func handleFsClone(p *Plugin, params json.RawMessage) (interface{}, error) {
	type fsCloneArgs struct {
		Fs    *FileSystem         `json:"src_fs"`
		Name  string              `json:"dest_fs_name"`
		Ss    *FileSystemSnapShot `json:"snapshot"`
		Flags uint64              `json:"flags"`
	}

	var args fsCloneArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_clone", uE)
	}

	fs, job, err := p.cb.File.FsClone(args.Fs, args.Name, args.Ss)
	return exclusiveOr(fs, job, err)

}

func handleFsFileClone(p *Plugin, params json.RawMessage) (interface{}, error) {
	type fsFileCloneArgs struct {
		Fs    *FileSystem         `json:"fs"`
		Src   string              `json:"src_file_name"`
		Dst   string              `json:"dest_file_name"`
		Ss    *FileSystemSnapShot `json:"snapshot"`
		Flags uint64              `json:"flags"`
	}

	var args fsFileCloneArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_clone", uE)
	}

	return p.cb.File.FsFileClone(args.Fs, args.Src, args.Dst, args.Ss)
}

func handleFsSnapShotCreate(p *Plugin, params json.RawMessage) (interface{}, error) {
	type fsSnapShotCreateArgs struct {
		Fs    *FileSystem `json:"fs"`
		Name  string      `json:"snapshot_name"`
		Flags uint64      `json:"flags"`
	}

	var args fsSnapShotCreateArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_snapshot_create", uE)
	}

	fs, job, err := p.cb.File.FsSnapShotCreate(args.Fs, args.Name)
	return exclusiveOr(fs, job, err)
}

func handleFsSnapShotDelete(p *Plugin, params json.RawMessage) (interface{}, error) {
	type fsSnapShotDeleteArgs struct {
		Fs    *FileSystem         `json:"fs"`
		Ss    *FileSystemSnapShot `json:"snapshot"`
		Flags uint64              `json:"flags"`
	}

	var args fsSnapShotDeleteArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_snapshot_delete", uE)
	}

	return p.cb.File.FsSnapShotDelete(args.Fs, args.Ss)
}

func handleFsSnapShots(p *Plugin, params json.RawMessage) (interface{}, error) {

	type fsSnapShotArgs struct {
		Fs    *FileSystem `json:"fs"`
		Flags uint64      `json:"flags"`
	}

	var args fsSnapShotArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_snapshots", uE)
	}

	return p.cb.File.FsSnapShots(args.Fs)
}

func handleFsSnapShotRestore(p *Plugin, params json.RawMessage) (interface{}, error) {

	type fsSnapShotRestoreArgs struct {
		Fs           *FileSystem         `json:"fs"`
		Ss           *FileSystemSnapShot `json:"snapshot"`
		All          bool                `json:"all_files"`
		Files        []string            `json:"files"`
		RestoreFiles []string            `json:"restore_files"`
		Flags        uint64              `json:"flags"`
	}

	var args fsSnapShotRestoreArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_snapshot_restore", uE)
	}

	return p.cb.File.FsSnapShotRestore(args.Fs, args.Ss, args.All, args.Files, args.RestoreFiles)
}

func handleFsHasChildDep(p *Plugin, params json.RawMessage) (interface{}, error) {
	type fsHasChildDepsArgs struct {
		Fs    *FileSystem `json:"fs"`
		Files []string    `json:"files"`
		Flags uint64      `json:"flags"`
	}

	var args fsHasChildDepsArgs
	if uE := json.Unmarshal(params, &args); uE != nil {
		return nil, invalidArgs("fs_child_dependency", uE)
	}

	return p.cb.File.FsHasChildDep(args.Fs, args.Files)
}

func nilAssign(present interface{}, cb handler) handler {

	// This seems like an epic fail of golang as I got burned by doing present == nil
	// ref. https://groups.google.com/forum/#!topic/golang-nuts/wnH302gBa4I/discussion
	if present == nil || reflect.ValueOf(present).IsNil() {
		return nil
	}
	return cb
}

func buildTable(c *PluginCallBacks) map[string]handler {
	return map[string]handler{
		"plugin_info":       handlePluginInfo,
		"plugin_register":   nilAssign(c.Mgmt.PluginRegister, handleRegister),
		"plugin_unregister": nilAssign(c.Mgmt.PluginUnregister, handleUnRegister),
		"systems":           nilAssign(c.Mgmt.Systems, handleSystems),
		"capabilities":      nilAssign(c.Mgmt.Capabilities, handleCapabilities),
		"time_out_set":      nilAssign(c.Mgmt.TimeOutSet, handleTmoSet),
		"time_out_get":      nilAssign(c.Mgmt.TimeOutGet, handleTmoGet),
		"pools":             nilAssign(c.Mgmt.Pools, handlePools),
		"job_status":        nilAssign(c.Mgmt.JobStatus, handleJobStatus),
		"job_free":          nilAssign(c.Mgmt.JobFree, handleJobFree),

		"volume_create":                      nilAssign(c.San.VolumeCreate, handleVolumeCreate),
		"volume_delete":                      nilAssign(c.San.VolumeDelete, handleVolumeDelete),
		"volumes":                            nilAssign(c.San.Volumes, handleVolumes),
		"disks":                              nilAssign(c.San.Disks, handleDisks),
		"volume_replicate":                   nilAssign(c.San.VolumeReplicate, handleVolumeReplicate),
		"volume_replicate_range":             nilAssign(c.San.VolumeReplicateRange, handleVolumeReplicateRange),
		"volume_replicate_range_block_size":  nilAssign(c.San.VolumeRepRangeBlkSize, handleVolRepRangeBlockSize),
		"volume_resize":                      nilAssign(c.San.VolumeResize, handleVolumeResize),
		"volume_enable":                      nilAssign(c.San.VolumeEnable, handleVolumeEnable),
		"volume_disable":                     nilAssign(c.San.VolumeDisable, handleVolumeDisable),
		"volume_mask":                        nilAssign(c.San.VolumeMask, handleVolumeMask),
		"volume_unmask":                      nilAssign(c.San.VolumeUnMask, handleVolumeUnMask),
		"volume_child_dependency":            nilAssign(c.San.VolHasChildDep, handleVolHasChildDep),
		"volume_child_dependency_rm":         nilAssign(c.San.VolChildDepRm, handleVolChildDepRm),
		"volumes_accessible_by_access_group": nilAssign(c.San.VolsMaskedToAg, handleVolsMaskedToAg),
		"access_groups":                      nilAssign(c.San.AccessGroups, handleAccessGroups),
		"access_group_create":                nilAssign(c.San.AccessGroupCreate, handleAccessGroupCreate),
		"access_group_delete":                nilAssign(c.San.AccessGroupDelete, handleAccessGroupDelete),
		"access_group_initiator_add":         nilAssign(c.San.AccessGroupInitAdd, handleAccessGroupInitAdd),
		"access_group_initiator_delete":      nilAssign(c.San.AccessGroupInitDelete, handleAccessGroupInitDelete),
		"access_groups_granted_to_volume":    nilAssign(c.San.AgsGrantedToVol, handleAgsGrantedToVol),
		"iscsi_chap_auth":                    nilAssign(c.San.IscsiChapAuthSet, handleIscsiChapAuthSet),
		"target_ports":                       nilAssign(c.San.TargetPorts, handleTargetPorts),

		"fs":                     nilAssign(c.File.FileSystems, handleFs),
		"fs_child_dependency":    nilAssign(c.File.FsHasChildDep, handleFsHasChildDep),
	}
}
