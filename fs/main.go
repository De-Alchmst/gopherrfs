package fs

import (
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
	Type fuse.DirentType

	Contents []fs.Node
}
type File struct{
	Inode uint64
	Name string
	Type fuse.DirentType

	Writer func([]byte) error
	Reader func() ([]byte, error)
}
// used for resolving
type path struct {
	FullPath string
}


var (
	rootDir = root{}
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
	return newPath(name), nil
}


func (root) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return []fuse.Dirent{}, nil
}
