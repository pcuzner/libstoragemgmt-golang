package localdisk

// #include <stdio.h>
// #include <libstoragemgmt/libstoragemgmt.h>
// #include <stdlib.h>
// #include <string.h>
import "C"
import (
	"unsafe"

	"github.com/libstorage/libstoragemgmt-golang/errors"
)

func processError(errorNum int, e *C.lsm_error) error {
	if e != nil {
		// Make sure we only free e if e is not nil
		defer C.lsm_error_free(e)
		return &errors.LsmError{
			Code:    int32(C.lsm_error_number_get(e)),
			Message: C.GoString(C.lsm_error_message_get(e))}
	}
	if errorNum != 0 {
		return &errors.LsmError{
			Code: int32(errorNum)}
	}
	return nil
}

func getStrings(lsmStrings *C.lsm_string_list, free bool) []string {
	var rc []string

	var num = C.lsm_string_list_size(lsmStrings)

	var i C.uint
	for i = 0; i < num; i++ {
		var item = C.GoString(C.lsm_string_list_elem_get(lsmStrings, i))
		rc = append(rc, item)
	}

	if free {
		C.lsm_string_list_free(lsmStrings)
	}
	return rc
}

// List returns local disk path(s)
func List() ([]string, error) {
	var disks []string

	var diskPaths *C.lsm_string_list
	var lsmError *C.lsm_error

	var e = C.lsm_local_disk_list(&diskPaths, &lsmError)
	if e == 0 {
		disks = getStrings(diskPaths, true)
	} else {
		return disks, processError(int(e), lsmError)
	}
	return disks, nil
}

// Vpd83Seach seaches local disks for vpd
func Vpd83Seach(vpd string) ([]string, error) {

	cs := C.CString(vpd)
	defer C.free(unsafe.Pointer(cs))

	var deviceList []string

	var slist *C.lsm_string_list
	var lsmError *C.lsm_error

	var err = C.lsm_local_disk_vpd83_search(cs, &slist, &lsmError)

	if err == 0 {
		deviceList = getStrings(slist, true)
	} else {
		return deviceList, processError(int(err), lsmError)
	}

	return deviceList, nil
}

// SerialNumGet retrieves the serial number for the local
// disk with the specfified path
func SerialNumGet(diskPath string) (string, error) {
	dp := C.CString(diskPath)
	defer C.free(unsafe.Pointer(dp))

	var sn *C.char
	var lsmError *C.lsm_error

	var rc = C.lsm_local_disk_serial_num_get(dp, &sn, &lsmError)
	if rc == 0 {
		var serialNum = C.GoString(sn)
		C.free(unsafe.Pointer(sn))
		return serialNum, nil
	}
	return "", processError(int(rc), lsmError)
}

