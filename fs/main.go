package fs

import (
	"fmt"
	"context"
	"os"

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

	defer c.Close()

	err = fs.Serve(c, FS{})
	if err != nil {
		return err
	}

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
	fmt.Println("Looking up:", name)
	handleEntry(name)
	return Path{FullPath: name}, nil
}


func (Root) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return []fuse.Dirent{}, nil
}
