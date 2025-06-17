package fs

import (
	"path/filepath"
	"context"
	"time"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"

	"gopherrfs/gopher"
)


type entry struct {
	Contents []byte
	TTL int64
	Status int
	Using int
}


const (
	entryStatusOK = iota
	entryStatusFailed
	entryStatusProcessing
)


var (
	DefaultTTL int64 = 60 * 60
	entries = map[string]*entry{}
)


func (p path) Attr(ctx context.Context, a *fuse.Attr) error {
	ent, ok := entries[p.FullPath]
	// File
	if ok { 
		ent.Using += 1

		for ent.Status == entryStatusProcessing {
			time.Sleep(10 * time.Millisecond)
		}

		a.Inode = 1
		a.Mode = os.ModeIrregular | 0o770
		a.Size = uint64(len(ent.Contents))

		ent.TTL = DefaultTTL
		ent.Using -= 1

	// Directory
	} else {
		a.Inode = 1
		a.Mode = os.ModeDir | 0o777
	}

	return nil
}


func (p path) ReadAll(ctx context.Context) ([]byte, error) {
	ent, ok := entries[p.FullPath]
	if !ok {
		return nil, fuse.ENOENT
	}

	ent.Using += 1
	for ent.Status == entryStatusProcessing {
		time.Sleep(10 * time.Millisecond)
	}

	ent.TTL = DefaultTTL
	ent.Using -= 1
	return ent.Contents, nil
}


func (p path) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	ent, ok := entries[p.FullPath]
	if !ok {
		return fuse.ENOENT
	}

	for ent.Status == entryStatusProcessing {
		time.Sleep(10 * time.Millisecond)
	}

	ent.Using += 1

	fuseutil.HandleRead(req, resp, ent.Contents)

	ent.TTL = DefaultTTL
	ent.Using -= 1

	return nil
}


func (path) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return []fuse.Dirent{}, nil
}


func (p path) Lookup(ctx context.Context, name string) (fs.Node, error) {
	return newPath(filepath.Join(p.FullPath, name)), nil
}


func newPath(name string) path {
	if name[len(name)-1] == ':' {
		handleEntry(name)
	}

	return path{FullPath: name}
}


func handleEntry(name string) {
	ent, ok := entries[name]

	if !ok {
		entries[name] = &entry{
			Contents: []byte{},
			TTL: DefaultTTL,
			Status: entryStatusProcessing,
			Using: 0,
		}

		go fillEntry(entries[name], name)

	} else {
		ent.TTL = DefaultTTL
	}
}


func fillEntry(ent *entry, name string) {
	data, err := gopher.FetchData(name[:len(name)-1]) // Remove trailing colon
	if err != nil {
		ent.Status = entryStatusFailed
	} else {
		ent.Status = entryStatusOK
	}

	ent.Contents = data
}
