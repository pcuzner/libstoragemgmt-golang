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
	class        string `json:"class"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	Status       uint32 `json:"status"`
	StatusInfo   string `json:"statis_info"`
	pluginData   string `json:"plugin_data"`
	FwVersion    string `json:"fw_version"`
	ReadCachePct int8   `json:"read_cache_pct"`
	SystemMode   int8   `json:"mode"`
}

// Volume represents a storage volume, aka. a logical unit
type Volume struct {
	class       string  `json:"class"`
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Enabled     LsmBool `json:"admin_state"`
	BlockSize   uint64  `json:"block_size"`
	NumOfBlocks uint64  `json:"num_of_blocks"`
	pluginData  string  `json:"plugin_data"`
	Vpd83       string  `json:"vpd83"`
	SystemID    string  `json:"system_id"`
	PoolID      string  `json:"pool_id"`
}

// Pool represents the unit of storage where block
// devices and/or file systems are created from.
type Pool struct {
	class              string `json:"class"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	ElementType        uint64 `json:"element_type"`
	UnsupportedActions uint64 `json:"unsupported_actions"`
	TotalSpace         uint64 `json:"total_space"`
	FreeSpace          uint64 `json:"free_space"`
	Status             uint64 `json:"status"`
	StatusInfo         string `json:"status_info"`
	pluginData         string `json:"plugin_data"`
	SystemID           string `json:"system_id"`
}

// Disk represents a physical device.
type Disk struct {
	class       string       `json:"class"`
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	DiskType    DiskType     `json:"disk_type"`
	BlockSize   uint64       `json:"block_size"`
	NumOfBlocks uint64       `json:"num_of_blocks"`
	Status      uint64       `json:"status"`
	pluginData  string       `json:"plugin_data"`
	SystemID    string       `json:"system_id"`
	Location    string       `json:"location"`
	Rpm         int          `json:"rpm"`
	LinkType    DiskLinkType `json:"link_type"`
	Vpd83       string       `json:"vpd83"`
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

// These constants are bitfields, eg. more than one bit can be set at the same time.
const (
	// DiskStatusUNKNOWN Plugin failed to query out the status of disk.
	DiskStatusUNKNOWN uint64 = 1

	// DiskStatusOk Disk is up and healthy.
	DiskStatusOk uint64 = 1 << 1

	//DiskStatusOther Vendor specific status.
	DiskStatusOther uint64 = 1 << 2

	//DiskStatusPredictiveFailure Disk is functional but will fail soon
	DiskStatusPredictiveFailure uint64 = 1 << 3

	//DiskStatusError Disk is not functional
	DiskStatusError uint64 = 1 << 4

	//DiskStatusRemoved Disk was removed by administrator
	DiskStatusRemoved uint64 = 1 << 5

	// DiskStatusStarting Disk is in the process of becomming ready.
	DiskStatusStarting uint64 = 1 << 6

	// DiskStatusStopping Disk is shutting down.
	DiskStatusStopping uint64 = 1 << 7

	// DiskStatusStopped Disk is stopped by administrator.
	DiskStatusStopped uint64 = 1 << 8

	// DiskStatusInitializing Disk is not yet functional, could be initializing eg. RAID, zeroed or scrubed etc.
	DiskStatusInitializing uint64 = 1 << 9

	// DiskStatusMaintenanceMode In maintenance for bad sector scan, integrity check and etc
	DiskStatusMaintenanceMode uint64 = 1 << 10

	// DiskStatusSpareDisk Disk is configured as a spare disk.
	DiskStatusSpareDisk uint64 = 1 << 11

	// DiskStatusReconstruct Disk is reconstructing its data.
	DiskStatusReconstruct uint64 = 1 << 12

	// DiskStatusFree Disk is not holding any data and it not designated as a spare.
	DiskStatusFree uint64 = 1 << 13
)
