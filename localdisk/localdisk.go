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
