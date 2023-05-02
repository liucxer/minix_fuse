package minix_decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
)

// 以下这些是st_mode 值的符号名称。
// 文件类型：
//#define S_IFMT 00170000		// 文件类型（8 进制表示）。
//#define S_IFREG 0100000		// 常规文件。
//#define S_IFBLK 0060000		// 块特殊（设备）文件，如磁盘dev/fd0。
//#define S_IFDIR 0040000		// 目录文件。
//#define S_IFCHR 0020000		// 字符设备文件。
//#define S_IFIFO 0010000		// FIFO 特殊文件。
//// 文件属性位：
//#define S_ISUID 0004000		// 执行时设置用户ID（set-user-ID）。
//#define S_ISGID 0002000		// 执行时设置组ID。
//#define S_ISVTX 0001000		// 对于目录，受限删除标志。
//
//#define S_ISREG(m) (((m) & S_IFMT) == S_IFREG)	// 测试是否常规文件。
//#define S_ISDIR(m) (((m) & S_IFMT) == S_IFDIR)	// 是否目录文件。
//#define S_ISCHR(m) (((m) & S_IFMT) == S_IFCHR)	// 是否字符设备文件。
//#define S_ISBLK(m) (((m) & S_IFMT) == S_IFBLK)	// 是否块设备文件。
//#define S_ISFIFO(m) (((m) & S_IFMT) == S_IFIFO)	// 是否FIFO 特殊文件。
//
//#define S_IRWXU 00700		// 宿主可以读、写、执行/搜索。
//#define S_IRUSR 00400		// 宿主读许可。
//#define S_IWUSR 00200		// 宿主写许可。
//#define S_IXUSR 00100		// 宿主执行/搜索许可。
//
//#define S_IRWXG 00070		// 组成员可以读、写、执行/搜索。
//#define S_IRGRP 00040		// 组成员读许可。
//#define S_IWGRP 00020		// 组成员写许可。
//#define S_IXGRP 00010		// 组成员执行/搜索许可。
//
//#define S_IRWXO 00007		// 其他人读、写、执行/搜索许可。
//#define S_IROTH 00004		// 其他人读许可。
//#define S_IWOTH 00002		// 其他人写许可。
//#define S_IXOTH 00001		// 其他人执行/搜索许可。

// Mode 总共16位   前面4位：文件类型    中间3位：文件属性位   后9位: 代表权限
type Mode uint16

type FileType int64

const (
	FILE_TYPE_IFREG  FileType = 8
	FILE_TYPE_IFBLK           = 6
	FILE_TYPE_IFDIR           = 4
	FILE_TYPE_IFCHAR          = 2
	FILE_TYPE_IFIFO           = 1
)

func (m Mode) FileType() FileType {
	//#define S_IFMT 00170000		// 文件类型（8 进制表示）。
	//#define S_IFREG 0100000		// 常规文件。
	//#define S_IFBLK 0060000		// 块特殊（设备）文件，如磁盘dev/fd0。
	//#define S_IFDIR 0040000		// 目录文件。
	//#define S_IFCHR 0020000		// 字符设备文件。
	//#define S_IFIFO 0010000		// FIFO 特殊文件。
	fileType := FileType(m >> 12)
	return fileType
}

func (m Mode) ISUID() bool {
	tmp := (m << 4) >> 9
	return tmp == 0x4
}

func (m Mode) ISGID() bool {
	tmp := (m << 4) >> 9
	return tmp == 0x2
}

func (m Mode) ISVTX() bool {
	tmp := (m << 4) >> 9
	return tmp == 0x1
}

func (m Mode) RGW() uint32 {
	tmp := m & 0x1FF
	return uint32(tmp)
}

// IsReg 是普通文件
func (m Mode) IsReg() bool {
	return m&0xF000 == 0x8000
}

// IsBlk 是块设备
func (m Mode) IsBlk() bool {
	return m&0xF000 == 0x6000
}

// IsDir 是目录
func (m Mode) IsDir() bool {
	return m&0xF000 == 0x4000
}

// IsChar 是字符设备
func (m Mode) IsChar() bool {
	return m&0xF000 == 0x2000
}

// IsFIFO 是管道文件
func (m Mode) IsFIFO() bool {
	return m&0xF000 == 0x1000
}

// Inode  683K大小
type Inode struct {
	Mode   Mode      // 文件类型和属性(rwx 位)。
	Uid    uint16    // 用户id（文件拥有者标识符）。
	Size   uint32    // 文件大小（字节数）。
	Time   uint32    // 修改时间（自1970.1.1:0 算起，秒）。
	Gid    uint8     // 组id(文件拥有者所在的组)。
	NLinks uint8     // 链接数（多少个文件目录项指向该i 节点）。
	Zone   [9]uint16 // 直接(0-6)、间接(7)或双重间接(8)逻辑块号。
}

type InodeItem struct {
	DirEntry []DirEntry
	Inode
	Data []byte `json:"data"`
}

func (inodeItem *InodeItem) String() string {
	var res string

	if inodeItem.Mode.IsDir() {
		res += "  inodeItem:"
		res += fmt.Sprintf("  inode: %+v", inodeItem.Inode)
		res += "   inode is dir. "
		res += fmt.Sprintf("  DirEntry: %+v", inodeItem.DirEntry)
		res += "\n"
	}

	if inodeItem.Mode.IsReg() {
		res += "  inodeItem:"
		res += fmt.Sprintf("  inode: %+v", inodeItem.Inode)
		res += "   inode is reg file. "
		res += fmt.Sprintf("  filedata: %s", string(inodeItem.Data))
	}

	return res
}

func (inodeItem *InodeItem) ReadData(bts []byte) error {
	var err error
	for i := 0; i < 7; i++ {
		if inodeItem.Zone[i] == 0 {
			break
		}
		start := int64(inodeItem.Zone[i]) * 1024
		end := int64(inodeItem.Zone[i])*1024 + 1024
		data := bts[start:end]
		inodeItem.Data = append(inodeItem.Data, data...)
	}

	if inodeItem.Zone[7] != 0 {
		start := inodeItem.Zone[7] * 1024
		end := inodeItem.Zone[7]*1024 + 1024
		inodeTable := bts[start:end]
		zoneList := make([]uint16, 1024/16)
		buf := bytes.NewBuffer(inodeTable)
		err := binary.Read(buf, binary.LittleEndian, &zoneList)
		if err != nil {
			logrus.Errorf("binary.Read err:%v", err)
			return err
		}

		for _, item := range zoneList {
			if item == 0 {
				break
			}
			start := item * 1024
			end := item*1024 + 1024
			inodeItem.Data = append(inodeItem.Data, bts[start:end]...)
		}
	}

	return err
}

type InodeTable struct {
	InodeItems []InodeItem `json:"inodeItems"`
}

func (inodeTable *InodeTable) String() string {
	var (
		res string
	)

	res += "inodeTable:\n"
	for i := 0; i < len(inodeTable.InodeItems); i++ {
		inodeItem := inodeTable.InodeItems[i]
		if inodeItem.Mode == 0 {
			continue
		}
		if inodeItem.Mode.IsDir() || inodeItem.Mode.IsReg() {
			res += fmt.Sprintf("    inodeNo:%d, %s", i, inodeItem.String())
		} else {
			res += fmt.Sprintf("    inodeNo:%d", i)
		}
	}

	return res
}

func (inodeTable *InodeTable) Decode(bts []byte, allBts []byte, inodeNum int64) error {
	var (
		err error
	)

	buf := bytes.NewBuffer(bts)
	inodeList := make([]Inode, inodeNum)
	err = binary.Read(buf, binary.LittleEndian, inodeList)
	if err != nil {
		logrus.Errorf("binary.Read err:%v", err)
		return err
	}
	for _, inode := range inodeList {
		var inodeItem InodeItem
		inodeItem.Inode = inode
		inodeTable.InodeItems = append(inodeTable.InodeItems, inodeItem)
	}

	for i := 0; i < len(inodeTable.InodeItems); i++ {
		err = (&inodeTable.InodeItems[i]).ReadData(allBts)
		if err != nil {
			logrus.Errorf("inode.ReadData err:%v", err)
			return err
		}

		if inodeTable.InodeItems[i].Mode.IsDir() {
			tmpBuf := bytes.NewBuffer(inodeTable.InodeItems[i].Data)
			dirEntry := make([]DirEntry, 32)
			err = binary.Read(tmpBuf, binary.LittleEndian, dirEntry)
			if err != nil {
				logrus.Errorf("binary.Read err:%v", err)
				return err
			}

			for _, dirEntryItem := range dirEntry {
				if dirEntryItem.InodeNo == 0 {
					continue
				}

				inodeTable.InodeItems[i].DirEntry = append(inodeTable.InodeItems[i].DirEntry, dirEntryItem)
			}
		}
	}

	return err
}
