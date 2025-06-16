package fs

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)


type FS struct{}
type Dir struct{}
type File struct{}


var (
	root = &Dir{}
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
