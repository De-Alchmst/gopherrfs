package fs

import (
	"fmt"
	"context"
	"os"
	"os/signal"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)


type filesystem struct{}
// used for static api
type root struct{}
type Dir struct{
	Inode uint64
	Name string

	Contents []DirNode
}
type File struct{
	Inode uint64
	Name string

	OnWrite func([]byte) error
	OnRead  func() ([]byte, error)
}
// used for resolving
type path struct {
	FullPath string
}

type fileHandle struct {
	Parent *File
	Inode uint64
	Contents []byte
}

type DirNode interface {
	fs.Node
	GetDirEnt() fuse.Dirent
	GetName() string
}


var (
	rootDir = root{}
	confDir  = &Dir{
		Inode: 0,
		Name: ":c",
		Contents: []DirNode{
			&File{
				Inode: 1,
				Name: "flush",
				OnWrite: func(data []byte) error {
					fmt.Println("Flushing entries...")
					return nil
				},
				OnRead: func() ([]byte, error) {
					return []byte("write stuff to flush cache\n"), nil
				},
			},
		},
	}
)


func MountFS(mountpoint, fsName, fsSubtype string) error {
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName(fsName),
		fuse.Subtype(fsSubtype),
	)

	if err != nil {
		return err
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Start serving in a goroutine
	serveDone := make(chan error, 1)
	go func() {
		serveDone <- fs.Serve(c, filesystem{})
	}()

	select {
		case <-serveDone:
		case <-sigChan:
	}
	
	err = fuse.Unmount(mountpoint)
	c.Close()

	return nil
}


func (filesystem) Root() (fs.Node, error) {
	return rootDir, nil
}


func (root) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0o555
	return nil
}


func (root) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == ":c" {
		return confDir, nil
	}
	return newPath(name), nil
}


func (root) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return []fuse.Dirent{confDir.GetDirEnt()}, nil
}
