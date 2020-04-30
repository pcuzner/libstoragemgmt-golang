Work in progress ...


```go

// A small sample to list systems
var c, _ = lsm.Client("sim://", "", 30000)
var systems, _ = c.Systems()

for _, s := range systems {
    fmt.Printf("%+v\n", s)
}
```