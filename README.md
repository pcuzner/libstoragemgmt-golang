[![Build Status](https://travis-ci.org/libstorage/libstoragemgmt-golang.svg?branch=master)](https://travis-ci.org/libstorage/libstoragemgmt-golang)

Work in progress, but almost complete from a client library perspective.


An example listing systems
```go
package main

import (
	"fmt"

	lsm "github.com/libstorage/libstoragemgmt-golang"
)

func main() {
	// Ignoring errors for succinctness
	c, _ := lsm.Client("sim://", "", 30000)
	systems, _ := c.Systems()
	for _, s := range systems {
		fmt.Printf("ID: %s, Name:%s, Version: %s\n", s.ID, s.Name, s.FwVersion)
	}
}

```

Example plugin library use can be found here: [https://github.com/tasleson/simgo](https://github.com/tasleson/simgo)