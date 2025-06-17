package main

import (
	"fmt"
	"flag"
	"log"
	"os"

	"gopherrfs/fs"
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

	err := fs.MountFS(mountpoint, "gopherrfs", "gopherrfs")

	if err != nil {
		log.Fatal(err)
	}
}
