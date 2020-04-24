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

// Disks represents a physical device.
type Disk struct {
	class       string `json:"class"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	DiskType    int    `json:"disk_type"`
	BlockSize   uint64 `json:"block_size"`
	NumOfBlocks uint64 `json:"num_of_blocks"`
	Status      uint64 `json:"status"`
	pluginData  string `json:"plugin_data"`
	SystemID    string `json:"system_id"`
	Location    string `json:"location"`
	Rpm         int    `json:"rpm"`
	LinkType    int    `json:"link_type"`
	Vpd83       string `json:"vpd83"`
}
