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

// JobFree instructs the plugin to release resources for the job that was returned.
func (c *ClientConnection) JobFree(jobID string) error {
	var args = make(map[string]interface{})
	args["job_id"] = jobID
	var result string
	var err = c.tp.invoke("job_free", args, &result)
	return err
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

// VolumeCreate creates a block device, returns job id, error.
// If job id and error are nil, then returnedVolume has newly created volume.
func (c *ClientConnection) VolumeCreate(
	pool *Pool,
	volumeName string,
	size uint64,
	provisioning VolumeProvisionType,
	sync bool,
	returnedVolume interface{}) (*string, error) {
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
	var invokeError = c.tp.invoke("capabilities", args, &cap)
	if invokeError != nil {
		return nil, invokeError
	}
	return &cap, nil
}

// VolumeResize resizes an existing volume, data loss may occur depending on storage implementation.
func (c *ClientConnection) VolumeResize(
	vol *Volume, newSizeBytes uint64, sync bool, returnedVolume interface{}) (*string, error) {
	var args = make(map[string]interface{})
	args["volume"] = *vol
	args["new_size_bytes"] = newSizeBytes

	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("volume_resize", args, &result), result, sync, returnedVolume)
}

// VolumeReplicate makes a replicated image of existing Volume
func (c *ClientConnection) VolumeReplicate(
	optionalPool *Pool, repType VolumeReplicateType, sourceVolume *Volume, name string,
	sync bool, returnedVolume interface{}) (*string, error) {

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

	var blkSize uint32
	var invokeError = c.tp.invoke("volume_replicate_range_block_size", args, &blkSize)
	if invokeError != nil {
		return 0, invokeError
	}
	return blkSize, nil
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
	returnedFs interface{}) (*string, error) {
	var args = make(map[string]interface{})
	args["pool"] = *pool
	args["name"] = name
	args["size_bytes"] = size

	var result [2]json.RawMessage
	return c.getJobOrResult(c.tp.invoke("fs_create", args, &result), result, sync, returnedFs)
}

// FsResize resizes an existing file system
func (c *ClientConnection) FsResize(
	fs *FileSystem, newSizeBytes uint64, sync bool, returnedFs interface{}) (*string, error) {
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
	returnedFs interface{}) (*string, error) {
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
func (c *ClientConnection) FsSnapShotCreate(fs *FileSystem, name string, sync bool, returnedSnapshot interface{}) (*string, error) {
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
