package fs

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)


type FS struct{}
// used for static api
type Root struct{}
type Dir struct{}
type File struct{}
// used for resolving
type Path struct {
	FullPath string
}


var (
	root = Root{}
)


func MountFS(mountpoint string) error {
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("gopherrfs"),
		fuse.Subtype("gopherrfs"),
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
		serveDone <- fs.Serve(c, FS{})
	}()

	select {
		case <-serveDone:
		case <-sigChan:
	}
	
	err = fuse.Unmount(mountpoint)
	c.Close()

	return nil
}


func (FS) Root() (fs.Node, error) {
	return root, nil
}


func (Root) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0o555
	return nil
}


func (Root) Lookup(ctx context.Context, name string) (fs.Node, error) {
	return newPath(name), nil
}


func (Root) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return []fuse.Dirent{}, nil
}
