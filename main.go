package main

import (
	"fmt"
	"flag"
	"log"
	"os"

	"gopherrfs/rfs"
)


func usage() {
	fmt.Println("Usage: gopherrfs <mountpoint>")
	flag.PrintDefaults()
}


func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}
	mountpoint := flag.Arg(0)

	err := rfs.MountFS(mountpoint, "gopherrfs", "gopherrfs",[]rfs.DirNode{}, API{})

	if err != nil {
		log.Fatal(err)
	}
}
