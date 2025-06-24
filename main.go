package main

import (
	"fmt"
	"flag"
	"log"
	"os"

	"github.com/de-alchmst/rfs"
)


func usage() {
	fmt.Println("Usage: gopherrfs <mountpoint>")
	flag.PrintDefaults()
}


func main() {
	// marse flags
	flag.Usage = usage
	var (
		defaultTTL = flag.Int("ttl", rfs.DefaultTTL, "TTL of cached entries")
		flushTime  = flag.Float64("flush", 5, "Time in seconds between TTL reduction")
	)
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	// set cache duration
	rfs.DefaultTTL = *defaultTTL
	rfs.CacheFlushTimeout = time.Duration(*flushTime * float64(time.Second))

	// activate the FS

	err := rfs.MountFS(mountpoint, "gopherrfs", "gopherrfs",[]rfs.DirNode{}, API{})

	if err != nil {
		log.Fatal(err)
	}
}
