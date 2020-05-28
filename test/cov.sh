#/bin/bash

LSM_GO_URI=simc:// go test -count 1 github.com/libstorage/libstoragemgmt-golang/test -coverpkg=../. -cover -coverprofile cc.out || exit 1
go tool cover -html=cc.out || exit 1

