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
	"time"

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

func GetAttrFromDecoderFile(file minix_decoder.File, a *fuse.Attr) error {
	var (
		err error
	)

	/*
		Blocks    uint64      // size in 512-byte units
		Rdev      uint32      // device numbers
		BlockSize uint32      // preferred blocksize for filesystem I/O
		Flags     AttrFlags
	*/
	a.Valid = time.Duration(0)
	a.Inode = uint64(file.InodeNo)
	a.Size = uint64(file.Inode.Size)
	a.Atime = time.Unix(int64(file.Inode.Time), 0)
	a.Mtime = time.Unix(int64(file.Inode.Time), 0)
	a.Ctime = time.Unix(int64(file.Inode.Time), 0)
	a.Nlink = uint32(file.Inode.NLinks)
	a.Uid = uint32(file.Inode.Uid)
	a.Gid = uint32(file.Inode.Gid)
	/*
		mode处理
	*/
	switch file.Inode.Mode.FileType() {
	case minix_decoder.FILE_TYPE_IFREG:
		a.Mode |= os.ModeType
	case minix_decoder.FILE_TYPE_IFBLK:
		a.Mode |= os.ModeDevice
	case minix_decoder.FILE_TYPE_IFDIR:
		a.Mode |= os.ModeDir
	case minix_decoder.FILE_TYPE_IFCHAR:
		a.Mode |= os.ModeCharDevice
	case minix_decoder.FILE_TYPE_IFIFO:
		a.Mode |= os.ModeNamedPipe
	}

	if file.Inode.Mode.ISUID() {
		a.Mode |= os.ModeSetuid
	}
	if file.Inode.Mode.ISGID() {
		a.Mode |= os.ModeSetgid
	}
	if file.Inode.Mode.ISVTX() {
		a.Mode |= os.ModeAppend
	}
	a.Mode |= os.FileMode(file.Inode.Mode.RGW())

	logrus.Infof("Attr:%+v", a)
	return err
}

func (dir Dir) Attr(ctx context.Context, attr *fuse.Attr) error {
	file, ok := FileMap[dir.InodeNo]
	if !ok {
		logrus.Errorf("inodeNo %d not exist", dir.InodeNo)
		return fmt.Errorf("inodeNo %d not exist", dir.InodeNo)
	}

	err := GetAttrFromDecoderFile(file, attr)
	if err != nil {
		logrus.Errorf("GetAttrFromDecoderFile err:%v", err)
		return err
	}

	return err
}

func (dir Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	var (
		node fs.Node
	)
	parentDir, ok := FileMap[dir.InodeNo]
	if !ok {
		logrus.Errorf("parentDir inodeNo %d not exist", dir.InodeNo)
		return node, fmt.Errorf("parentDir inodeNo %d not exist", dir.InodeNo)
	}

	subInodeNo := int64(0)
	for _, file := range parentDir.Files {
		if name == file.Path {
			subInodeNo = file.InodeNo
			break
		}
	}

	return File{
		InodeNo: subInodeNo,
		Name:    name,
	}, syscall.ENOENT
}

func (dir Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var (
		err     error
		dirents []fuse.Dirent
	)

	file, ok := FileMap[dir.InodeNo]
	if !ok {
		logrus.Errorf("inodeNo %d not exist", dir.InodeNo)
		return dirents, fmt.Errorf("inodeNo %d not exist", dir.InodeNo)
	}
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
type File struct {
	InodeNo int64
	Name    string
}

const greeting = "hello, world\n"

func (file File) Attr(ctx context.Context, a *fuse.Attr) error {
	decoderFile, ok := FileMap[file.InodeNo]
	if !ok {
		logrus.Errorf("inodeNo %d not exist", file.InodeNo)
		return fmt.Errorf("inodeNo %d not exist", file.InodeNo)
	}

	err := GetAttrFromDecoderFile(decoderFile, a)
	if err != nil {
		logrus.Errorf("GetAttrFromDecoderFile err:%v", err)
		return err
	}

	return err
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
