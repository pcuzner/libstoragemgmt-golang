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
	PluginData   string `json:"plugin_data"`
	FwVersion    string `json:"fw_version"`
	ReadCachePct int8   `json:"read_cache_pct"`
	SystemMode   int8   `json:"mode"`
}
