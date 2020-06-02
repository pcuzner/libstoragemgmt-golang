#/bin/bash

if [[ -z "${LSMSRC}" ]]; then
    export CGO_LDFLAGS="-lstoragemgmt"
else
    export CGO_CFLAGS="-I$LSMSRC/c_binding/include"
    export CGO_LDFLAGS="-L$LSMSRC/c_binding/.libs/ -lstoragemgmt"
    export LD_LIBRARY_PATH="$LSMSRC/c_binding/.libs/"
fi

export LSM_GO_URI="simc://"

go test -count 1 github.com/libstorage/libstoragemgmt-golang/test -coverpkg=../. -cover -coverprofile client.out || exit 1
go test -count 1 github.com/libstorage/libstoragemgmt-golang/test -coverpkg=.././localdisk -cover -coverprofile localdisk.out || exit 1

go tool cover -html=client.out || exit 1
go tool cover -html=localdisk.out || exit 1

exit 0
