package fs

import (
	"fmt"
	"path/filepath"
	"context"
	"time"
	"os"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	"gopherrfs/gopher"
)


type Entry struct {
	Contents []byte
	Ttl int64
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
	fmt.Println("Resolving path:", p.FullPath)
	entry := Entries[p.FullPath]
	entry.Using += 1

	for entry.Status == EntryStatusProcessing {
		time.Sleep(10 * time.Millisecond)
	}

	a.Inode = 1
	a.Mode = os.ModeIrregular | 0o774
	a.Size = uint64(len(entry.Contents))

	entry.Ttl = DefaultTTL
	entry.Using -= 1
	return nil
}


func (p Path) Lookup(ctx context.Context, name string) (fs.Node, error) {
	fmt.Println("Looking up:", name)
	newPath := filepath.Join(p.FullPath, name)
	handleEntry(newPath)

	return Path{FullPath: newPath}, nil
}


func (Path) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	return []fuse.Dirent{}, nil
}


func handleEntry(path string) {
	entry, ok := Entries[path]

	if !ok {
		Entries[path] = &Entry{
			Contents: []byte{},
			Ttl: DefaultTTL,
			Status: EntryStatusProcessing,
			Using: 0,
		}

		go fillEntry(Entries[path], path)

	} else {
		entry.Ttl = DefaultTTL
	}
}


func fillEntry(entry *Entry, path string) {
	data, err := gopher.FetchData(path)
	if err != nil {
		entry.Status = EntryStatusFailed
		return
	}

	entry.Contents = data
	entry.Status = EntryStatusOK
}
