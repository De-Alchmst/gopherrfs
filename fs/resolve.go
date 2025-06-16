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


type Entry struct {
	Contents []byte
	TTL int64
	Status int
	Using int
}


const (
	EntryStatusOK = iota
	EntryStatusFailed
	EntryStatusProcessing
)

const (
	DefaultTTL = 60 * 60
)


var (
	Entries = map[string]*Entry{}
)


func (p Path) Attr(ctx context.Context, a *fuse.Attr) error {
	entry, ok := Entries[p.FullPath]
	// File
	if ok { 
		entry.Using += 1

		for entry.Status == EntryStatusProcessing {
			time.Sleep(10 * time.Millisecond)
		}

		a.Inode = 1
		a.Mode = os.ModeIrregular | 0o770
		a.Size = uint64(len(entry.Contents))

		entry.TTL = DefaultTTL
		entry.Using -= 1

	// Directory
	} else {
		a.Inode = 1
		a.Mode = os.ModeDir | 0o770
	}

	return nil
}


func (p Path) ReadAll(ctx context.Context) ([]byte, error) {
	entry, ok := Entries[p.FullPath]
	if !ok {
		return nil, fuse.ENOENT
	}

	entry.Using += 1
	for entry.Status == EntryStatusProcessing {
		time.Sleep(10 * time.Millisecond)
	}

	entry.TTL = DefaultTTL
	entry.Using -= 1
	return entry.Contents, nil
}


func (p Path) Read(ctx context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	entry, ok := Entries[p.FullPath]
	if !ok {
		return fuse.ENOENT
	}

	for entry.Status == EntryStatusProcessing {
		time.Sleep(10 * time.Millisecond)
	}

	entry.Using += 1

	fuseutil.HandleRead(req, resp, entry.Contents)

	entry.TTL = DefaultTTL
	entry.Using -= 1

	return nil
}


func (Path) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return []fuse.Dirent{}, nil
}


func (p Path) Lookup(ctx context.Context, name string) (fs.Node, error) {
	return newPath(filepath.Join(p.FullPath, name)), nil
}


func newPath(path string) Path {
	if path[len(path)-1] == ':' {
		handleEntry(path)
	}

	return Path{FullPath: path}
}


func handleEntry(path string) {
	entry, ok := Entries[path]

	if !ok {
		Entries[path] = &Entry{
			Contents: []byte{},
			TTL: DefaultTTL,
			Status: EntryStatusProcessing,
			Using: 0,
		}

		go fillEntry(Entries[path], path)

	} else {
		entry.TTL = DefaultTTL
	}
}


func fillEntry(entry *Entry, path string) {
	data, err := gopher.FetchData(path[:len(path)-1]) // Remove trailing colon
	if err != nil {
		entry.Status = EntryStatusFailed
		return
	}

	entry.Contents = data
	entry.Status = EntryStatusOK
}
