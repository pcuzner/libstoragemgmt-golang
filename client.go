// SPDX-License-Identifier: 0BSD

package libstoragemgmt

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	errors "github.com/libstorage/libstoragemgmt-golang/errors"
)

// ClientConnection is the structure that encomposes the needed data for the plugin connection.
type ClientConnection struct {
	tp         transPort
	PluginName string
	timeout    uint32
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

	return &ClientConnection{tp: *transport, PluginName: pluginName, timeout: timeout}, nil
}

// PluginInfo information about the current plugin
func (c *ClientConnection) PluginInfo() (*PluginInfo, error) {
	var args = make(map[string]interface{})
	var info []string
	var invokeError = c.tp.invoke("plugin_info", args, &info)
	if invokeError != nil {
		return nil, invokeError
	}
	return &PluginInfo{Description: info[0], Version: info[1], Name: c.PluginName}, nil
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
	return systems, c.tp.invoke("systems", args, &systems)
}

// Volumes returns block device information
func (c *ClientConnection) Volumes(search ...string) ([]Volume, error) {
	var args = make(map[string]interface{})
	var volumes []Volume

	if !handleSearch(args, search) {
		return make([]Volume, 0), &errors.LsmError{
			Code: errors.InvalidArgument,
			Message: fmt.Sprintf(
				"volume supports 0 or 2 search parameters (key, value), provide %d", len(search)),
			Data: ""}
	}

	return volumes, c.tp.invoke("volumes", args, &volumes)
}

// Pools returns the units of storage that block devices and FS
// can be created from.
func (c *ClientConnection) Pools(search ...string) ([]Pool, error) {
	var args = make(map[string]interface{})

	if !handleSearch(args, search) {
		return make([]Pool, 0), &errors.LsmError{
			Code: errors.InvalidArgument,
			Message: fmt.Sprintf(
				"pool supports 0 or 2 search parameters (key, value), provide %d", len(search)),
			Data: ""}
	}

	var pools []Pool
	return pools, c.tp.invoke("pools", args, &pools)
}

// Disks returns disks that are present.
func (c *ClientConnection) Disks() ([]Disk, error) {
	var args = make(map[string]interface{})
	var disks []Disk
	return disks, c.tp.invoke("disks", args, &disks)
}

// FileSystems returns pools that are present.
func (c *ClientConnection) FileSystems() ([]FileSystem, error) {
	var args = make(map[string]interface{})
	var fileSystems []FileSystem
	return fileSystems, c.tp.invoke("fs", args, &fileSystems)
}

// NfsExports returns nfs exports  that are present.
func (c *ClientConnection) NfsExports(search ...string) ([]NfsExport, error) {
	var args = make(map[string]interface{})

	if !handleSearch(args, search) {
		return make([]NfsExport, 0), &errors.LsmError{
			Code: errors.InvalidArgument,
			Message: fmt.Sprintf(
				"NfsExports supports 0 or 2 search parameters (key, value), provide %d", len(search)),
			Data: ""}
	}

	var nfsExports []NfsExport
	return nfsExports, c.tp.invoke("exports", args, &nfsExports)
}

// NfsExportAuthTypes returns list of support authentication types
func (c *ClientConnection) NfsExportAuthTypes() ([]string, error) {
	var args = make(map[string]interface{})
	var authTypes []string
	var err = c.tp.invoke("export_auth", args, &authTypes)
	return authTypes, err
}

// FsExport creates or modifies a NFS export.
func (c *ClientConnection) FsExport(fs *FileSystem, exportPath *string,
	access *NfsAccess, authType *string, options *string, nfsExport *NfsExport) error {

	if len(access.Ro) == 0 && len(access.Rw) == 0 {
		return &errors.LsmError{
			Code:    errors.InvalidArgument,
			Message: "at least 1 host should exist in access.ro or access.rw",
			Data:    ""}
	}

	for _, i := range access.Root {
		if !contains(access.Rw, i) && !contains(access.Ro, i) {
			return &errors.LsmError{
				Code: errors.InvalidArgument,
				Message: fmt.Sprintf(
					"host '%s' contained in access.root should also be in access.rw or access.ro", i),
				Data: ""}
		}
	}

	for _, i := range access.Rw {
		if contains(access.Ro, i) {
			return &errors.LsmError{
				Code:    errors.InvalidArgument,
				Message: fmt.Sprintf("host '%s' in both access.ro and access.rw", i),
				Data:    ""}
		}
	}

	var args = make(map[string]interface{})
	args["fs_id"] = fs.ID
	args["export_path"] = exportPath
	args["root_list"] = emptySliceIfNil(access.Ro)
	args["rw_list"] = emptySliceIfNil(access.Rw)
	args["ro_list"] = emptySliceIfNil(access.Ro)
	args["anon_uid"] = access.AnonUID
	args["anon_gid"] = access.AnonGID
	args["auth_type"] = authType
	args["options"] = options

	return c.tp.invoke("export_fs", args, nfsExport)
}

// FsUnExport removes a file system export.
func (c *ClientConnection) FsUnExport(export *NfsExport) error {
	var args = make(map[string]interface{})
	args["export"] = *export
	return c.tp.invoke("export_remove", args, nil)
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
	return targetPorts, c.tp.invoke("target_ports", args, &targetPorts)
}

// Batteries returns batteries that are present
func (c *ClientConnection) Batteries() ([]Battery, error) {
	var args = make(map[string]interface{})
	var batteries []Battery
	return batteries, c.tp.invoke("batteries", args, &batteries)
}

// JobFree instructs the plugin to release resources for the job that was returned.
func (c *ClientConnection) JobFree(jobID string) error {
	var args = make(map[string]interface{})
	args["job_id"] = jobID
	var result string
	return c.tp.invoke("job_free", args, &result)
}

// JobStatus instructs the plugin to return the status of the specified job.  The returned values are
// the current job status, percent complete, and any errors that occured.  Always check error first as if it's
// set the other two are meaningless.  If checking on the status of an operation that doesn't return a result
// or you are not wanting the result, pass nil.
func (c *ClientConnection) JobStatus(jobID string, returnedResult interface{}) (JobStatusType, uint8, error) {
	var args = make(map[string]interface{})
	args["job_id"] = jobID

	var result [3]json.RawMessage
	var jobError = c.tp.invoke("job_status", args, &result)
	if jobError != nil {
		return JobStatusError, 0, jobError
	}

	var status JobStatusType
	var statusMe = json.Unmarshal(result[0], &status)
	if statusMe != nil {
		return JobStatusError, 0, statusMe
	}

	switch status {
	case JobStatusInprogress:
		var percent uint8
		var percentError = json.Unmarshal(result[1], &percent)
		if percentError != nil {
			return JobStatusError, 0, percentError
		}
		return status, percent, nil
	case JobStatusComplete:
		// Some RPC calls with jobs do not return a value, thus the third item is
		// "null"
		if string(result[2]) != "null" && returnedResult != nil {
			return status, 100, json.Unmarshal(result[2], returnedResult)
		}
		return status, 100, nil
	case JobStatusError:
		// Error
		var error errors.LsmError
		var checkErrorE = json.Unmarshal(result[2], &error)
		if checkErrorE != nil {
			return JobStatusError, 0, checkErrorE
		}
		return JobStatusError, 0, &errors.LsmError{
			Code:    errors.PluginBug,
			Message: "job_status returned error status with no error information"}
	default:
		return JobStatusError, 0, &errors.LsmError{
			Code:    errors.PluginBug,
			Message: fmt.Sprintf("Invalid status type returned %v", status)}
	}
}

func (c *ClientConnection) getJobOrResult(err error, returned [2]json.RawMessage, sync bool, result interface{}) (*string, error) {
	if err != nil {
		return nil, err
	}

	var job string
	var um = json.Unmarshal(returned[0], &job)
	if um == nil {
		// We have a job, but want to wait for result, so do so.
		if sync {
			return nil, c.JobWait(job, result)
		}

		return &job, nil
	}
	// We have the result
	var umO = json.Unmarshal(returned[1], result)
	return nil, umO
}

func (c *ClientConnection) getJobOrNone(err error, returned json.RawMessage, sync bool) (*string, error) {
	if err != nil {
		return nil, err
	}

	var job string
	var um = json.Unmarshal(returned, &job)
	if um == nil {
		// We have a job, but want to wait for result, so do so.
		if sync {
			return nil, c.JobWait(job, nil)
		}

		return &job, nil
	}
	return nil, um
}

// JobWait waits for the job to finish and retrieves the end result in "returnedResult".
func (c *ClientConnection) JobWait(jobID string, returnedResult interface{}) error {

	for true {
		var status, _, err = c.JobStatus(jobID, returnedResult)
		if err != nil {
			return err
		}

		if status == JobStatusInprogress {
			time.Sleep(time.Millisecond * 250)
			continue
		} else if status == JobStatusComplete {
			var freeError = c.JobFree(jobID)
			if freeError != nil {
				return &errors.LsmError{
					Code: errors.PluginBug,
					Message: fmt.Sprintf(
						"We successfully waited for job %s, but got an error freeing it: %s", jobID, freeError)}
			}
			break
		}
	}
	return nil
}

// TimeOutSet sets the connection timeout with the storage device.
func (c *ClientConnection) TimeOutSet(milliSeconds uint32) error {
	var args = make(map[string]interface{})
	args["ms"] = milliSeconds
	var err = c.tp.invoke("time_out_set", args, nil)
	if err == nil {
		c.timeout = milliSeconds
	}
	return err
}

// TimeOutGet sets the connection timeout with the storage device.
func (c *ClientConnection) TimeOutGet() uint32 {
	return c.timeout
}

// VolumeCreate creates a block device, returns job id, error.
// If job id and error are nil, then returnedVolume has newly created volume.
func (c *ClientConnection) VolumeCreate(
	pool *Pool,
	volumeName string,
	size uint64,
	provisioning VolumeProvisionType,
	sync bool,
	returnedVolume *Volume) (*string, error) {
	var args = make(map[string]interface{})
	args["pool"] = *pool
	args["volume_name"] = volumeName
	args["size_bytes"] = size
	args["provisioning"] = provisioning

	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("volume_create", args, &result), result, sync, returnedVolume)
}

// VolumeDelete deletes a block device.
func (c *ClientConnection) VolumeDelete(vol *Volume, sync bool) (*string, error) {
	var args = make(map[string]interface{})
	args["volume"] = *vol
	var result json.RawMessage
	return c.getJobOrNone(c.tp.invoke("volume_delete", args, &result), result, sync)
}

// Capabilities retrieve capabilities
func (c *ClientConnection) Capabilities(system *System) (*Capabilities, error) {
	var args = make(map[string]interface{})
	args["system"] = *system
	var cap Capabilities
	return &cap, c.tp.invoke("capabilities", args, &cap)
}

// VolumeResize resizes an existing volume, data loss may occur depending on storage implementation.
func (c *ClientConnection) VolumeResize(
	vol *Volume, newSizeBytes uint64, sync bool, returnedVolume *Volume) (*string, error) {
	var args = make(map[string]interface{})
	args["volume"] = *vol
	args["new_size_bytes"] = newSizeBytes

	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("volume_resize", args, &result), result, sync, returnedVolume)
}

// VolumeReplicate makes a replicated image of existing Volume
func (c *ClientConnection) VolumeReplicate(
	optionalPool *Pool, repType VolumeReplicateType, sourceVolume *Volume, name string,
	sync bool, returnedVolume *Volume) (*string, error) {

	var args = make(map[string]interface{})
	if optionalPool != nil {
		args["pool"] = *optionalPool
	} else {
		args["pool"] = nil
	}
	args["volume_src"] = *sourceVolume
	args["rep_type"] = repType
	args["name"] = name

	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("volume_replicate", args, &result), result, sync, returnedVolume)
}

// VolumeRepRangeBlkSize block size for replicating a range of blocks
func (c *ClientConnection) VolumeRepRangeBlkSize(system *System) (uint32, error) {
	var args = make(map[string]interface{})
	args["system"] = *system

	var blkSize uint32 = 0
	return blkSize, c.tp.invoke("volume_replicate_range_block_size", args, &blkSize)
}

// VolumeReplicateRange replicates a range of blocks on the same or different Volume
func (c *ClientConnection) VolumeReplicateRange(
	repType VolumeReplicateType, srcVol *Volume, dstVol *Volume,
	ranges []BlockRange, sync bool) (*string, error) {

	var args = make(map[string]interface{})
	args["rep_type"] = repType
	args["ranges"] = ranges
	args["volume_src"] = *srcVol
	args["volume_dest"] = *dstVol

	var result json.RawMessage
	return c.getJobOrNone(c.tp.invoke("volume_replicate_range", args, &result), result, sync)
}

// FsCreate creates a file system, returns job id, error.
// If job id and error are nil, then returnedFs has newly created filesystem.
func (c *ClientConnection) FsCreate(
	pool *Pool,
	name string,
	size uint64,
	sync bool,
	returnedFs *FileSystem) (*string, error) {
	var args = make(map[string]interface{})
	args["pool"] = *pool
	args["name"] = name
	args["size_bytes"] = size

	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("fs_create", args, &result), result, sync, returnedFs)
}

// FsResize resizes an existing file system
func (c *ClientConnection) FsResize(
	fs *FileSystem, newSizeBytes uint64, sync bool, returnedFs *FileSystem) (*string, error) {
	var args = make(map[string]interface{})
	args["fs"] = *fs
	args["new_size_bytes"] = newSizeBytes

	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("fs_resize", args, &result), result, sync, returnedFs)
}

// FsDelete deletes a file system.
func (c *ClientConnection) FsDelete(fs *FileSystem, sync bool) (*string, error) {
	var args = make(map[string]interface{})
	args["fs"] = *fs
	var result json.RawMessage
	return c.getJobOrNone(c.tp.invoke("fs_delete", args, &result), result, sync)
}

// FsClone makes a clone of an existing file system
func (c *ClientConnection) FsClone(
	srcFs *FileSystem,
	destName string,
	optionalSnapShot *FileSystemSnapShot,
	sync bool,
	returnedFs *FileSystem) (*string, error) {
	var args = make(map[string]interface{})
	args["src_fs"] = *srcFs
	args["dest_fs_name"] = destName

	if optionalSnapShot != nil {
		args["snapshot"] = *optionalSnapShot
	} else {
		args["snapshot"] = nil
	}

	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("fs_clone", args, &result), result, sync, returnedFs)
}

// FsFileClone makes a clone of an existing file system
func (c *ClientConnection) FsFileClone(
	fs *FileSystem,
	srcFileName string,
	dstFileName string,
	optionalSnapShot *FileSystemSnapShot,
	sync bool,
) (*string, error) {
	var args = make(map[string]interface{})

	args["fs"] = *fs
	args["src_file_name"] = srcFileName
	args["dest_file_name"] = dstFileName

	if optionalSnapShot != nil {
		args["snapshot"] = *optionalSnapShot
	} else {
		args["snapshot"] = nil
	}

	var result json.RawMessage
	return c.getJobOrNone(c.tp.invoke("fs_file_clone", args, &result), result, sync)
}

// FsSnapShotCreate creates a file system snapshot for the supplied snapshot
// If job id and error are nil, then returnedFs has newly created filesystem.
func (c *ClientConnection) FsSnapShotCreate(fs *FileSystem, name string, sync bool,
	returnedSnapshot *FileSystemSnapShot) (*string, error) {
	var args = make(map[string]interface{})
	args["fs"] = *fs
	args["snapshot_name"] = name
	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("fs_snapshot_create", args, &result), result, sync, returnedSnapshot)
}

// FsSnapShotDelete deletes a file system snapshot.
func (c *ClientConnection) FsSnapShotDelete(fs *FileSystem, snapShot *FileSystemSnapShot, sync bool) (*string, error) {
	var args = make(map[string]interface{})
	args["fs"] = *fs
	args["snapshot"] = *snapShot
	var result json.RawMessage
	return c.getJobOrNone(c.tp.invoke("fs_snapshot_delete", args, &result), result, sync)
}

// FsSnapShots returns list of file system snapsthos for specified file system.
// can be created from.
func (c *ClientConnection) FsSnapShots(fs *FileSystem) ([]FileSystemSnapShot, error) {
	var args = make(map[string]interface{})
	var snapShots []FileSystemSnapShot

	args["fs"] = *fs

	var err = c.tp.invoke("fs_snapshots", args, &snapShots)
	return snapShots, err
}

// FsSnapShotRestore restores all the files for a file systems or specific files.
func (c *ClientConnection) FsSnapShotRestore(
	fs *FileSystem, snapShot *FileSystemSnapShot, allFiles bool,
	files []string, restoreFiles []string, sync bool) (*string, error) {

	var args = make(map[string]interface{})

	if !allFiles {
		if len(files) == 0 {
			return nil, &errors.LsmError{
				Code:    errors.InvalidArgument,
				Message: "'files' is empty and 'all_files' is false!"}
		}

		if len(files) != len(restoreFiles) {
			return nil, &errors.LsmError{
				Code:    errors.InvalidArgument,
				Message: "'files' and 'restoreFiles' have different lengths!"}
		}
	}

	args["fs"] = *fs
	args["snapshot"] = *snapShot
	args["files"] = files
	args["restore_files"] = restoreFiles
	args["all_files"] = allFiles
	var result json.RawMessage
	return c.getJobOrNone(c.tp.invoke("fs_snapshot_restore", args, &result), result, sync)
}

// FsHasChildDep checks whether file system has a child dependency.
func (c *ClientConnection) FsHasChildDep(fs *FileSystem, files []string) (bool, error) {

	var args = make(map[string]interface{})

	args["fs"] = *fs
	args["files"] = files
	var result bool
	var err = c.tp.invoke("fs_child_dependency", args, &result)
	return result, err
}

// FsChildDepRm remove dependencies for specified file system.
func (c *ClientConnection) FsChildDepRm(
	fs *FileSystem, files []string, sync bool) (*string, error) {

	var args = make(map[string]interface{})
	args["fs"] = *fs
	args["files"] = files
	var result json.RawMessage
	return c.getJobOrNone(c.tp.invoke("fs_child_dependency_rm", args, &result), result, sync)
}

// AccessGroupCreate creates an access group.
func (c *ClientConnection) AccessGroupCreate(name string, initID string,
	initType InitiatorType, system *System, accessGroup *AccessGroup) error {

	var args = make(map[string]interface{})

	var check = validateInitID(initID, initType)
	if check != nil {
		return check
	}

	args["name"] = name
	args["init_id"] = initID
	args["init_type"] = initType
	args["system"] = *system
	return c.tp.invoke("access_group_create", args, accessGroup)
}

// AccessGroupDelete deletes an access group.
func (c *ClientConnection) AccessGroupDelete(ag *AccessGroup) error {
	var args = make(map[string]interface{})
	args["access_group"] = *ag
	var result json.RawMessage
	return c.tp.invoke("access_group_delete", args, result)
}

// AccessGroupInitAdd adds an initiator to an access group.
func (c *ClientConnection) AccessGroupInitAdd(ag *AccessGroup,
	initID string, initType InitiatorType, accessGroup *AccessGroup) error {

	var check = validateInitID(initID, initType)
	if check != nil {
		return check
	}

	var args = make(map[string]interface{})
	args["access_group"] = *ag
	args["init_id"] = initID
	args["init_type"] = initType
	return c.tp.invoke("access_group_initiator_add", args, accessGroup)
}

// AccessGroupInitDelete deletes an initiator from an access group.
func (c *ClientConnection) AccessGroupInitDelete(ag *AccessGroup,
	initID string, initType InitiatorType, accessGroup *AccessGroup) error {

	var check = validateInitID(initID, initType)
	if check != nil {
		return check
	}

	var args = make(map[string]interface{})
	args["access_group"] = *ag
	args["init_id"] = initID
	args["init_type"] = initType
	return c.tp.invoke("access_group_initiator_delete", args, accessGroup)
}
