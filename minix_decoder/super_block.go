package minix_decoder

import (
	"bytes"
	"encoding/binary"
	"github.com/sirupsen/logrus"
)

// SuperBlock 超级块 1K大小
type SuperBlock struct {
	InodeNum             uint16 /* total number of inodes */        // 0x5560  = 21856 inodes
	ZoneNum              uint16 /* total number of zones */         // 0xffff  = 65535 blocks
	InodeBitmapBlocksNum uint16 /* number of inode bitmap blocks */ // 0x0003  inode位图占用3个block       0x1400 - 0x0800 = 0xC00 = 3 * 1024
	ZoneBitmapBlocks     uint16 /* number of zone bitmap blocks */  // 0x0008  data block位图占用8个block  0x3400 - 0x1400 = 0x2000 = 8 * 1024
	FirstDataZone        uint16 /* first data zone */               // 0x02b8  = 696
	ZoneSize             uint16 // 0x0000                            // 一个block占用2的0次方K，代表1k
	MaxFileSize          uint32 /* maximum file size */ // 0x10081c00 = 268966912
	Magic                uint16 // 0x138f 魔幻数字
	State                uint16 // 是否挂在，0: 已经挂载， 1: 未挂载
	Zones                uint32
}

func (superBlock *SuperBlock) Decode(bts []byte) error {
	var (
		err error
	)
	buf := bytes.NewBuffer(bts)
	err = binary.Read(buf, binary.LittleEndian, superBlock)
	if err != nil {
		logrus.Errorf("binary.Read err:%v", err)
		return err
	}
	return err
}
