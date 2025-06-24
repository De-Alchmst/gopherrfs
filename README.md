# GopherRFS

GopherRFS is one of the reference
[resolving file systems](https://github.com/De-Alchmst/resolving-file-system-spec/)
using the
[Go rfs module](https://pkg.go.dev/github.com/de-alchmst/rfs).

Other reference file systems include:

* [WebRFS](https://github.com/De-Alchmst/webrfs)

## Build

```
go mod tidy
go build
```

## use

```
Usage: gopherrfs [flags] <mountpoint>
  -flush float
        Time in seconds between TTL reduction (default 5)
  -ttl int
        TTL of cached entries (default 60)
```

GopherRFS does not include any extra modifiers and does not support writing.
