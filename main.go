// Hellofs implements a simple "hello world" file system.
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/liucxer/minix_fuse/minix_decoder"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"syscall"

	"github.com/liucxer/minix_fuse/fuse"
	"github.com/liucxer/minix_fuse/fuse/fs"
	_ "github.com/liucxer/minix_fuse/fuse/fs/fstestutil"
)

// FS implements the hello world file system.
type FS struct{}

func (FS) Root() (fs.Node, error) {
	return Dir{
		InodeNo: 0,
	}, nil
}

// Dir implements both Node and Handle for the root directory.
type Dir struct {
	InodeNo int64
}

func (Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	file := FileMap[0]
	a.Inode = uint64(file.InodeNo)
	a.Mode = os.FileMode(file.Inode.Mode)
	return nil
}

func (Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if name == "hello" {
		return File{}, nil
	}
	return nil, syscall.ENOENT
}

func (dir Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var (
		err     error
		dirents []fuse.Dirent
	)

	file := FileMap[dir.InodeNo]
	for _, subFile := range file.Files {
		var dirent fuse.Dirent
		dirent.Inode = uint64(subFile.InodeNo)
		dirent.Name = subFile.Path
		subfile := FileMap[subFile.InodeNo]
		if subfile.Inode.Mode.IsDir() {
			dirent.Type = fuse.DT_Dir
		} else if subfile.Inode.Mode.IsReg() {
			dirent.Type = fuse.DT_File
		}
		dirents = append(dirents, dirent)
	}
	return dirents, err
}

// File implements both Node and Handle for the hello file.
type File struct{}

const greeting = "hello, world\n"

func (File) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 2
	a.Mode = 0o444
	a.Size = uint64(len(greeting))
	return nil
}

func (File) ReadAll(ctx context.Context) ([]byte, error) {
	return []byte(greeting), nil
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "  %s MOUNTPOINT\n", os.Args[0])
	flag.PrintDefaults()
}

var FileMap map[int64]minix_decoder.File

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		os.Exit(2)
	}

	mountPoint := flag.Arg(0)
	err := fuse.Unmount(mountPoint)
	if err != nil {
		logrus.Errorf("fuse.Unmount err:%v", err)
	}

	c, err := fuse.Mount(
		mountPoint,
		fuse.FSName("minix_fuse"),
		fuse.Subtype("minix_fuse"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	devicePath := flag.Arg(1)
	fileMap, err := minix_decoder.GetFileMap(devicePath)
	if err != nil {
		logrus.Errorf("minix_decoder.GetFiles err:%v", err)
		return
	}
	FileMap = fileMap

	err = fs.Serve(c, FS{})
	if err != nil {
		logrus.Errorf("fs.Serve err:%v", err)
		return
	}
}
