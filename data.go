// SPDX-License-Identifier: 0BSD

package libstoragemgmt

// PluginInfo - Information about a specific plugin
type PluginInfo struct {
	Version     string
	Description string
	Name        string
}

// System represents a storage system.
// * A hardware RAID card
// * A storage area network (SAN)
// * A software solutions running on commidity hardware
// * A Linux system running NFS Service
type System struct {
	Class        string           `json:"class"`
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Status       SystemStatusType `json:"status"`
	StatusInfo   string           `json:"status_info"`
	PluginData   string           `json:"plugin_data"`
	FwVersion    string           `json:"fw_version"`
	ReadCachePct int8             `json:"read_cache_pct"`
	SystemMode   SystemModeType   `json:"mode"`
}

// SystemModeType type representing system mode.
type SystemModeType int8

// SystemStatusType type representing system status.
type SystemStatusType uint32

const (
	// SystemReadCachePctNoSupport System read cache percentage not supported.
	SystemReadCachePctNoSupport int8 = -2

	// SystemReadCachePctUnknown System read cache percentage unknown.
	SystemReadCachePctUnknown int8 = -1

	// SystemStatusUnknown System status is unknown.
	SystemStatusUnknown SystemStatusType = 1

	// SystemStatusOk  System status is OK.
	SystemStatusOk SystemStatusType = 1 << 1

	// SystemStatusError System is in error state.
	SystemStatusError SystemStatusType = 1 << 2

	// SystemStatusDegraded System is degraded in some way
	SystemStatusDegraded SystemStatusType = 1 << 3

	// SystemStatusPredictiveFailure System has potential failure.
	SystemStatusPredictiveFailure SystemStatusType = 1 << 4

	// SystemStatusOther Vendor specific status.
	SystemStatusOther SystemStatusType = 1 << 5

	// SystemModeUnknown Plugin failed to query system mode.
	SystemModeUnknown SystemModeType = -2

	// SystemModeNoSupport Plugin does not support querying system mode.
	SystemModeNoSupport SystemModeType = -1

	//SystemModeHardwareRaid The storage system is a hardware RAID card
	SystemModeHardwareRaid SystemModeType = 0

	// SystemModeHba The physical disks can be exposed to OS directly without any
	// configurations.
	SystemModeHba SystemModeType = 1
)

// Volume represents a storage volume, aka. a logical unit
type Volume struct {
	Class       string  `json:"class"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Enabled     LsmBool `json:"admin_state"`
	BlockSize   uint64  `json:"block_size"`
	NumOfBlocks uint64  `json:"num_of_blocks"`
	PluginData  string  `json:"plugin_data"`
	Vpd83       string  `json:"vpd83"`
	SystemID    string  `json:"system_id"`
	PoolID      string  `json:"pool_id"`
}

// JobStatusType is enumerated type returned from Job control
type JobStatusType uint32

const (

	// JobStatusInprogress indicated job is in progress
	JobStatusInprogress JobStatusType = 1

	// JobStatusComplete indicates job is complete
	JobStatusComplete JobStatusType = 2

	// JobStatusError indicated job has errored
	JobStatusError JobStatusType = 3
)

// VolumeProvisionType enumerated type for volume creation provisioning
type VolumeProvisionType int

const (
	// VolumeProvisionTypeUnknown provision type unknown
	VolumeProvisionTypeUnknown VolumeProvisionType = -1

	// VolumeProvisionTypeThin thin provision volume
	VolumeProvisionTypeThin VolumeProvisionType = 1

	// VolumeProvisionTypeFull fully provision volume
	VolumeProvisionTypeFull VolumeProvisionType = 2

	// VolumeProvisionTypeDefault use the default for the storage provider
	VolumeProvisionTypeDefault VolumeProvisionType = 3
)

// Pool represents the unit of storage where block
// devices and/or file systems are created from.
type Pool struct {
	Class              string              `json:"class"`
	ID                 string              `json:"id"`
	Name               string              `json:"name"`
	ElementType        PoolElementType     `json:"element_type"`
	UnsupportedActions PoolUnsupportedType `json:"unsupported_actions"`
	TotalSpace         uint64              `json:"total_space"`
	FreeSpace          uint64              `json:"free_space"`
	Status             PoolStatusType      `json:"status"`
	StatusInfo         string              `json:"status_info"`
	PluginData         string              `json:"plugin_data"`
	SystemID           string              `json:"system_id"`
}

// PoolElementType type used to describe what things can be created from pool
type PoolElementType uint64

// PoolUnsupportedType type used to describe what actions are unsupported
type PoolUnsupportedType uint64

// PoolStatusType type used to describe the status of pool
type PoolStatusType uint64

const (

	// PoolElementPool This pool could allocate space for sub pool.
	PoolElementPool PoolElementType = 1 << 1

	// PoolElementTypeVolume This pool can be used for volume creation.
	PoolElementTypeVolume PoolElementType = 1 << 2

	// PoolElementTypeFs this pool can be used to for FS creation.
	PoolElementTypeFs PoolElementType = 1 << 3

	// PoolElementTypeDelta this pool can hold delta data for snapshots.
	PoolElementTypeDelta PoolElementType = 1 << 4

	// PoolElementTypeVolumeFull this pool could be used to create fully allocated volume.
	PoolElementTypeVolumeFull PoolElementType = 1 << 5

	// PoolElementTypeVolumeThin this pool could be used to create thin provisioned volume.
	PoolElementTypeVolumeThin PoolElementType = 1 << 6

	// PoolElementTypeSysReserved this pool is reserved for internal use.
	PoolElementTypeSysReserved PoolElementType = 1 << 10

	// PoolUnsupportedVolumeGrow this pool does not allow growing volumes
	PoolUnsupportedVolumeGrow PoolUnsupportedType = 1

	// PoolUnsupportedVolumeShink this pool does not allow shrinking volumes
	PoolUnsupportedVolumeShink PoolUnsupportedType = 1 << 1

	// PoolStatusUnknown Plugin failed to query pool status.
	PoolStatusUnknown PoolStatusType = 1

	// PoolStatusOk The data of this pool is accessible with no data loss. But it might
	// be set with PoolStatusDegraded to indicate redundancy loss.
	PoolStatusOk PoolStatusType = 1 << 1

	// PoolStatusOther Vendor specific status, check Pool.StatusInfo for more information.
	PoolStatusOther PoolStatusType = 1 << 2

	// PoolStatusDegraded indicates pool has lost data redundancy.
	PoolStatusDegraded PoolStatusType = 1 << 4

	// PoolStatusError indicates pool data is not accessible due to some members offline.
	PoolStatusError PoolStatusType = 1 << 5

	// PoolStatusStopped pool is stopped by administrator.
	PoolStatusStopped PoolStatusType = 1 << 9

	// PoolStatusStarting is reviving from STOPPED status. Pool data is not accessible yet.
	PoolStatusStarting PoolStatusType = 1 << 10

	// PoolStatusReconstructing pool is be reconstructing hash or mirror data.
	PoolStatusReconstructing PoolStatusType = 1 << 12

	// PoolStatusVerifying indicates array is running integrity check(s).
	PoolStatusVerifying PoolStatusType = 1 << 13

	// PoolStatusInitializing indicates pool is not accessible and performing initialization.
	PoolStatusInitializing PoolStatusType = 1 << 14

	// PoolStatusGrowing indicates pool is growing in size.  PoolStatusInfo can contain more
	// information about this task.  If PoolStatusOk is set, data is still accessible.
	PoolStatusGrowing PoolStatusType = 1 << 15
)

// Disk represents a physical device.
type Disk struct {
	class       string         `json:"class"`
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	DiskType    DiskType       `json:"disk_type"`
	BlockSize   uint64         `json:"block_size"`
	NumOfBlocks uint64         `json:"num_of_blocks"`
	Status      DiskStatusType `json:"status"`
	pluginData  string         `json:"plugin_data"`
	SystemID    string         `json:"system_id"`
	Location    string         `json:"location"`
	Rpm         int            `json:"rpm"`
	LinkType    DiskLinkType   `json:"link_type"`
	Vpd83       string         `json:"vpd83"`
}

// DiskType is an enumerated type representing different types of disks.
type DiskType int

const (
	// DiskTypeUnknown Plugin failed to query disk type
	DiskTypeUnknown DiskType = 0

	// DiskTypeOther Vendor specific disk type
	DiskTypeOther DiskType = 1

	// DiskTypeAta IDE disk type.
	DiskTypeAta DiskType = 3

	// DiskTypeSata SATA disk
	DiskTypeSata DiskType = 4

	// DiskTypeSas SAS disk.
	DiskTypeSas DiskType = 5

	// DiskTypeFc FC disk.
	DiskTypeFc DiskType = 6

	// DiskTypeSop SCSI over PCI-Express.
	DiskTypeSop DiskType = 7

	// DiskTypeScsi SCSI disk.
	DiskTypeScsi DiskType = 8

	// DiskTypeLun Remote LUN from SAN array.
	DiskTypeLun DiskType = 9

	// DiskTypeNlSas Near-Line SAS, just SATA disk + SAS port
	DiskTypeNlSas DiskType = 51

	// DiskTypeHdd Normal HDD, fall back value if failed to detect HDD type(SAS/SATA/etc).
	DiskTypeHdd DiskType = 52

	// DiskTypeSsd Solid State Drive.
	DiskTypeSsd DiskType = 53

	// DiskTypeHybrid Hybrid disk uses a combination of HDD and SSD.
	DiskTypeHybrid DiskType = 54
)

// DiskLinkType is an enumerated type representing different types of disk connection.
type DiskLinkType int

const (
	// DiskLinkTypeNoSupport Plugin does not support querying disk link type.
	DiskLinkTypeNoSupport DiskLinkType = -2

	// DiskLinkTypeUnknown Plugin failed to query disk link type
	DiskLinkTypeUnknown DiskLinkType = -1

	// DiskLinkTypeFc Fibre channel
	DiskLinkTypeFc DiskLinkType = 0

	//DiskLinkTypeSsa Serial Storage Architecture
	DiskLinkTypeSsa DiskLinkType = 2

	// DiskLinkTypeSbp Serial Bus Protocol, used by IEEE 1394.
	DiskLinkTypeSbp = 3

	// DiskLinkTypeSrp SCSI RDMA Protocol
	DiskLinkTypeSrp = 4

	// DiskLinkTypeIscsi Internet Small Computer System Interface
	DiskLinkTypeIscsi = 5

	// DiskLinkTypeSas Serial Attached SCSI.
	DiskLinkTypeSas = 6

	// DiskLinkTypeAdt Automation/Drive Interface Transport. Often used by tape.
	DiskLinkTypeAdt = 7

	// DiskLinkTypeAta PATA/IDE or SATA.
	DiskLinkTypeAta = 8

	// DiskLinkTypeUsb USB
	DiskLinkTypeUsb = 9

	// DiskLinkTypeSop SCSI over PCI-E.
	DiskLinkTypeSop = 10

	// DiskLinkTypePciE PCI-E, e.g. NVMe.
	DiskLinkTypePciE = 11
)

// DiskStatusType base type for bitfield
type DiskStatusType uint64

// These constants are bitfields, eg. more than one bit can be set at the same time.
const (
	// DiskStatusUNKNOWN Plugin failed to query out the status of disk.
	DiskStatusUNKNOWN DiskStatusType = 1

	// DiskStatusOk Disk is up and healthy.
	DiskStatusOk DiskStatusType = 1 << 1

	//DiskStatusOther Vendor specific status.
	DiskStatusOther DiskStatusType = 1 << 2

	//DiskStatusPredictiveFailure Disk is functional but will fail soon
	DiskStatusPredictiveFailure DiskStatusType = 1 << 3

	//DiskStatusError Disk is not functional
	DiskStatusError DiskStatusType = 1 << 4

	//DiskStatusRemoved Disk was removed by administrator
	DiskStatusRemoved DiskStatusType = 1 << 5

	// DiskStatusStarting Disk is in the process of becomming ready.
	DiskStatusStarting DiskStatusType = 1 << 6

	// DiskStatusStopping Disk is shutting down.
	DiskStatusStopping DiskStatusType = 1 << 7

	// DiskStatusStopped Disk is stopped by administrator.
	DiskStatusStopped DiskStatusType = 1 << 8

	// DiskStatusInitializing Disk is not yet functional, could be initializing eg. RAID, zeroed or scrubed etc.
	DiskStatusInitializing DiskStatusType = 1 << 9

	// DiskStatusMaintenanceMode In maintenance for bad sector scan, integrity check and etc
	DiskStatusMaintenanceMode DiskStatusType = 1 << 10

	// DiskStatusSpareDisk Disk is configured as a spare disk.
	DiskStatusSpareDisk DiskStatusType = 1 << 11

	// DiskStatusReconstruct Disk is reconstructing its data.
	DiskStatusReconstruct DiskStatusType = 1 << 12

	// DiskStatusFree Disk is not holding any data and it not designated as a spare.
	DiskStatusFree DiskStatusType = 1 << 13
)

// FileSystem represents a file systems information
type FileSystem struct {
	class      string `json:"class"`
	ID         string `json:"id"`
	Name       string `json:"name"`
	TotalSpace uint64 `json:"total_space"`
	FreeSpace  uint64 `json:"free_space"`
	pluginData string `json:plugin_data"`
	SystemID   string `json:"system_id"`
	PoolID     string `json:"pool_id"`
}

// NfsExport represents exported file systems over NFS.
type NfsExport struct {
	class       string   `json:"class"`
	ID          string   `json:"id"`
	FsID        string   `json:"fs_id"`
	ExportPath  string   `json:"export_path"`
	Auth        string   `json:"auth"`
	Root        []string `json:"root"`
	Rw          []string `json:"rw"`
	Ro          []string `json:"ro"`
	AnonUID     int64    `json:"anonuid"`
	Options     string   `json:"options"`
	plugin_data string   `json:"plugin_data"`
}

// AccessGroup represents a collection of initiators.
type AccessGroup struct {
	class         string        `json:"class"`
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	InitIDs       []string      `json:"init_ids"`
	InitiatorType InitiatorType `json:"init_type"`
	pluginData    string        `json:plugin_data"`
	SystemID      string        `json:"system_id"`
}

// InitiatorType is enumerated type of initiators
type InitiatorType int

const (
	// InitiatorTypeUnknown plugin failed to query initiator type
	InitiatorTypeUnknown InitiatorType = 0

	// InitiatorTypeOther vendor specific initiator type
	InitiatorTypeOther InitiatorType = 1

	// InitiatorTypeWwpn FC or FCoE WWPN
	InitiatorTypeWwpn InitiatorType = 2

	// InitiatorTypeIscsiIqn iSCSI IQN
	InitiatorTypeIscsiIqn InitiatorType = 5

	// InitiatorTypeMixed this access group contains more than 1 type of initiator
	InitiatorTypeMixed InitiatorType = 7
)

// TargetPort represents information about target ports.
type TargetPort struct {
	class           string   `json:"class"`
	ID              string   `json:"id"`
	PortType        PortType `json:"port_type"`
	ServiceAddress  string   `json:"service_address"`
	NetworkAddress  string   `json:"network_address"`
	PhysicalAddress string   `json:"physical_address"`
	PhysicalName    string   `json:"physical_name"`
	pluginData      string   `json:plugin_data"`
	SystemID        string   `json:"system_id"`
}

// PortType in enumerated type of port
type PortType int32

const (

	// PortTypeOther is a vendor specific port type
	PortTypeOther PortType = 1

	// PortTypeFc indicates FC port type
	PortTypeFc PortType = 2

	// PortTypeFCoE indicates FC over Ethernet type
	PortTypeFCoE PortType = 3

	// PortTypeIscsi indicates FC over iSCSI type
	PortTypeIscsi PortType = 4
)

// Battery represents a battery in the system.
type Battery struct {
	class       string        `json:"class"`
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	BatteryType BatteryType   `json:"type"`
	pluginData  string        `json:plugin_data"`
	Status      BatteryStatus `json:"status"`
	SystemID    string        `json:"system_id"`
}

// BatteryType indicates enumerated type of battery
type BatteryType int32

// BatteryStatus indicates bitfield for status of battery
type BatteryStatus uint64

const (
	// BatteryTypeUnknown plugin failed to detect battery type
	BatteryTypeUnknown BatteryType = 1

	// BatteryTypeOther vendor specific battery type
	BatteryTypeOther BatteryType = 2

	// BatteryTypeChemical indicates li-ion etc.
	BatteryTypeChemical BatteryType = 3

	// BatteryTypeCapacitor indicates capacitor
	BatteryTypeCapacitor BatteryType = 4

	// BatteryStatusUnknown plugin failed to query battery status
	BatteryStatusUnknown BatteryStatus = 1

	// BatteryStatusOther vendor specific status
	BatteryStatusOther BatteryStatus = 1 << 1

	// BatteryStatusOk indicated battery is healthy and operational
	BatteryStatusOk BatteryStatus = 1 << 2

	// BatteryStatusDischarging indicates battery is discharging
	BatteryStatusDischarging BatteryStatus = 1 << 3

	// BatteryStatusCharging battery is charging
	BatteryStatusCharging BatteryStatus = 1 << 4

	// BatteryStatusLearning indicated battery system is optimizing battery use
	BatteryStatusLearning BatteryStatus = 1 << 5

	// BatteryStatusDegraded indicated battery should be checked and/or replaced
	BatteryStatusDegraded BatteryStatus = 1 << 6

	// BatteryStatusError indicates battery is in bad state
	BatteryStatusError = 1 << 7
)

// Capabilities representation
type Capabilities struct {
	class string `json:"class"`
	Cap   string `json:"cap"`
}

// IsSupported used to determine if a capability is supported
func (c *Capabilities) IsSupported(cap CapabilityType) bool {
	var capIdx = int32(cap) * 2
	if c.Cap[capIdx:capIdx+2] == "01" {
		return true
	}
	return false
}

// CapabilityType Enumerated type for capabilities
type CapabilityType uint32

const (
	// CapVolumes supports retrieving Volumes
	CapVolumes CapabilityType = 20

	// CapVolumeCreate supports VolumeCreate
	CapVolumeCreate CapabilityType = 21

	// CapVolumeCResize supports VolumeResize
	CapVolumeCResize CapabilityType = 22

	// CapVolumeCReplicate supports VolumeReplicate
	CapVolumeCReplicate CapabilityType = 23

	// CapVolumeCReplicateClone supports volume replication via clone
	CapVolumeCReplicateClone CapabilityType = 24

	// CapVolumeCReplicateCopy supports volume replication via copy
	CapVolumeCReplicateCopy CapabilityType = 25

	// CapVolumeCReplicateMirrorAsync supports volume replication via async. mirror
	CapVolumeCReplicateMirrorAsync CapabilityType = 26

	// CapVolumeCReplicateMirrorSync supports volume replication via sync. mirror
	CapVolumeCReplicateMirrorSync CapabilityType = 27

	// CapVolumeCopyRangeBlockSize supports reporting of what block size to be used in Copy Range
	CapVolumeCopyRangeBlockSize CapabilityType = 28

	// CapVolumeCopyRange supports copying a range of a Volume
	CapVolumeCopyRange CapabilityType = 29

	// CapVolumeCopyRangeClone supports a range clone of a volume
	CapVolumeCopyRangeClone CapabilityType = 30

	// CapVolumeCopyRangeCopy supports a range copy of a volume
	CapVolumeCopyRangeCopy CapabilityType = 31

	// CapVolumeDelete supports volume deletion
	CapVolumeDelete CapabilityType = 33

	// CapVolumeEnable admin. volume enable
	CapVolumeEnable CapabilityType = 34

	// CapVolumeDisable admin. volume disable
	CapVolumeDisable CapabilityType = 35

	// CapVolumeMask support volume masking
	CapVolumeMask CapabilityType = 36

	// CapVolumeUnmask support volume unmasking
	CapVolumeUnmask CapabilityType = 37

	// CapAccessGroups supports AccessGroups
	CapAccessGroups CapabilityType = 38

	// CapAccessGroupCreateWwpn supports access group wwpn creation
	CapAccessGroupCreateWwpn CapabilityType = 39

	// CapAccessGroupDelete delete an access group
	CapAccessGroupDelete CapabilityType = 40

	// CapAccessGroupInitiatorAddWwpn support adding WWPN to an access group
	CapAccessGroupInitiatorAddWwpn CapabilityType = 41

	// CapAccessGroupInitiatorDel supports removal of an initiator from access group
	CapAccessGroupInitiatorDel CapabilityType = 42

	// CapVolumesMaskedToAg supports listing of volumes masked to access group
	CapVolumesMaskedToAg CapabilityType = 43

	// CapAgsGrantedToVol list access groups with access to volume
	CapAgsGrantedToVol CapabilityType = 44

	// CapHasChildDep indicates support for determing if volume has child dep.
	CapHasChildDep CapabilityType = 45

	// CapChildDepRm indiates support for removing child dep.
	CapChildDepRm CapabilityType = 46

	// CapAccessGroupCreateIscsiIqn supports ag creating with iSCSI initiator
	CapAccessGroupCreateIscsiIqn CapabilityType = 47

	// CapAccessGroupInitAddIscsiIqn supports adding iSCSI initiator
	CapAccessGroupInitAddIscsiIqn CapabilityType = 48

	// CapIscsiChapAuthSet support iSCSI CHAP setting
	CapIscsiChapAuthSet CapabilityType = 53

	// CapVolRaidInfo supports RAID information
	CapVolRaidInfo CapabilityType = 54

	// CapVolumeThin supports creating thinly provisioned Volumes.
	CapVolumeThin CapabilityType = 55

	// CapBatteries supports Batteries Call
	CapBatteries CapabilityType = 56

	// CapVolCacheInfo supports CacheInfo
	CapVolCacheInfo CapabilityType = 57

	// CapVolPhyDiskCacheSet support VolPhyDiskCacheSet
	CapVolPhyDiskCacheSet CapabilityType = 58

	// CapVolPhysicalDiskCacheSetSystemLevel supports VolPhyDiskCacheSet
	CapVolPhysicalDiskCacheSetSystemLevel CapabilityType = 59

	// CapVolWriteCacheSetEnable supports VolWriteCacheSet
	CapVolWriteCacheSetEnable CapabilityType = 60

	// CapVolWriteCacheSetAuto supports VolWriteCacheSet
	CapVolWriteCacheSetAuto CapabilityType = 61

	// CapVolWriteCacheSetDisabled supported VolWriteCacheSet
	CapVolWriteCacheSetDisabled CapabilityType = 62

	// CapVolWriteCacheSetImpactRead indicates the VolWriteCacheSet might also
	// impact read cache policy.
	CapVolWriteCacheSetImpactRead CapabilityType = 63

	// CapVolWriteCacheSetWbImpactOther Indicate the VolWriteCacheSet with
	// `wcp=Cache::Enabled` might impact other volumes in the same
	// system.
	CapVolWriteCacheSetWbImpactOther CapabilityType = 64

	// CapVolReadCacheSet Support VolReadCacheSet()
	CapVolReadCacheSet CapabilityType = 65

	// VolReadCacheSetImpactWrite Indicates the VolReadCacheSet might
	// also impact write cache policy.
	VolReadCacheSetImpactWrite CapabilityType = 66

	// CapFs support Fs listing.
	CapFs CapabilityType = 100

	// CapFsDelete supports  FsDelete
	CapFsDelete CapabilityType = 101

	// CapFsResize Support FsResize
	CapFsResize CapabilityType = 102

	// CapFsCreate support FsCreate
	CapFsCreate CapabilityType = 103

	// CapFsClone support FsClone
	CapFsClone CapabilityType = 104

	// CapFsFileClone support FsFileClone
	CapFsFileClone CapabilityType = 105

	// CapFsSnapshots support FsSnapshots
	CapFsSnapshots CapabilityType = 106

	// CapFsSnapshotCreate support FsSnapshotCreate
	CapFsSnapshotCreate CapabilityType = 107

	// CapFsSnapshotDelete support FfsSnapshotDelete
	CapFsSnapshotDelete CapabilityType = 109

	// CapFsSnapshotRestore support FsSnapshotRestore
	CapFsSnapshotRestore CapabilityType = 110

	// CapFsSnapshotRestoreSpecificFiles support FsSnapshotRestore with `files` argument.
	CapFsSnapshotRestoreSpecificFiles CapabilityType = 111

	// CapFsHasChildDep support FsHasChildDep
	CapFsHasChildDep CapabilityType = 112

	// CapFsChildDepRm support FsChildDepRm
	CapFsChildDepRm CapabilityType = 113

	// CapFsChildDepRmSpecificFiles support FsChildDepRm with `files` argument.
	CapFsChildDepRmSpecificFiles CapabilityType = 114

	// CapNfsExportAuthTypeList support NfsExpAuthTypeList
	CapNfsExportAuthTypeList CapabilityType = 120

	// CapNfsExports support NfsExports
	CapNfsExports CapabilityType = 121

	// CapFsExport support FsExport
	CapFsExport CapabilityType = 122

	// CapFsUnexport support FsUnexport
	CapFsUnexport CapabilityType = 123

	// CapFsExportCustomPath support FsExport
	CapFsExportCustomPath CapabilityType = 124

	// CapSysReadCachePctSet support SystemReadCachePctSet
	CapSysReadCachePctSet CapabilityType = 158

	// CapSysReadCachePctGet support Systems() `ReadCachePct` property
	CapSysReadCachePctGet CapabilityType = 159

	// CapSysFwVersionGet support Systems()  with valid `FwVersion` property.
	CapSysFwVersionGet CapabilityType = 160

	// CapSysModeGet support `Systems()` with valid `mode` property.
	CapSysModeGet CapabilityType = 161

	// CapDiskLocation support Disks with valid `Location` property.
	CapDiskLocation CapabilityType = 163

	// CapDiskRpm support `Disks()` with valid `Rpm` property.
	CapDiskRpm CapabilityType = 164

	// CapDiskLinkType support `Disks()` with valid `LinkType` property.
	CapDiskLinkType CapabilityType = 165

	// CapVolumeLed support `VolIdentLedOn()` and `VolIdentLedOff()`.
	CapVolumeLed CapabilityType = 171

	// CapTargetPorts Support TargetPorts()
	CapTargetPorts CapabilityType = 216

	// CapDisks support Disks()
	CapDisks CapabilityType = 220

	// CapPoolMemberInfo support `PoolMemberInfo()`
	CapPoolMemberInfo CapabilityType = 221

	// CapVolumeRaidCreate support `VolRaidCreateCapGet()` and
	// VolRaidCreate().
	CapVolumeRaidCreate CapabilityType = 222

	//CapDiskVpd83Get support `Disks()` with valid `Vpd83` property.
	CapDiskVpd83Get CapabilityType = 223
)
