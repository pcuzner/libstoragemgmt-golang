[![Build Status](https://travis-ci.org/libstorage/libstoragemgmt-golang.svg?branch=master)](https://travis-ci.org/libstorage/libstoragemgmt-golang)

Work in progress, but almost complete ...


An example listing systems
```go
package main

import (
	"fmt"

	lsm "github.com/libstorage/libstoragemgmt-golang"
)

func main() {
	// Ignoring errors for succinctness
	var c, _ = lsm.Client("sim://", "", 30000)
	var systems, _ = c.Systems()
	for _, s := range systems {
		fmt.Printf("ID: %s, Name:%s, Version: %s\n", s.ID, s.Name, s.FwVersion)
	}
}

```