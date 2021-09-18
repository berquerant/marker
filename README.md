# marker

Given

```go
package example

type Kern struct {}
```

run `marker -type Kern -method IsKern` then generate

```go
package example

func (*Kern) IsKern() {}
```

in Kern_marker.go in the same directory.
