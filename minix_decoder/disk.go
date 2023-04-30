package minix_decoder

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type DiskData struct {
	BootBlock       BootBlock       `json:"bootBlock"`       // 启动块 1K
	SuperBlock      SuperBlock      `json:"superBlock"`      // 超级块 1K
	InodeBitMap     InodeBitMap     `json:"inodeBitMap"`     // InodeBitMap  3K大小
	DataBlockBitMap DataBlockBitMap `json:"dataBlockBitMap"` // DataBlockBitMap  8K大小
	InodeTable      InodeTable      `json:"inodeTable"`      // Inode 表  683K大小
	DataBlock       DataBlock       `json:"dataBlock"`       // DataBlock  65535K大小
}

func (disk *DiskData) BootBlockData(buf []byte) []byte {
	return buf[0:1024]
}

func (disk *DiskData) SuperBlockData(buf []byte) []byte {
	return buf[1024:2048]
}

func (disk *DiskData) InodeBitMapData(buf []byte) []byte {
	return buf[2048 : 2048+disk.SuperBlock.InodeBitmapBlocksNum*1024]
}

func (disk *DiskData) DataBlockBitMapData(buf []byte) []byte {
	inodeBitMapOffset := 2048 + disk.SuperBlock.InodeBitmapBlocksNum*1024
	return buf[inodeBitMapOffset : inodeBitMapOffset+disk.SuperBlock.ZoneBitmapBlocks*1024]
}

func (disk *DiskData) InodeTableData(buf []byte) []byte {
	inodeTableOffset := int64(2048 + disk.SuperBlock.InodeBitmapBlocksNum*1024 + disk.SuperBlock.ZoneBitmapBlocks*1024)
	sizeofInode := int64(32)
	return buf[inodeTableOffset : inodeTableOffset+int64(disk.SuperBlock.InodeNum)*sizeofInode]
}

func (disk *DiskData) Decode(filePath string) error {
	var (
		err error
	)

	bts, err := ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("ioutil.ReadFile err:%v, path:%s", err, filePath)
		return err
	}

	// 解析超级块
	err = disk.SuperBlock.Decode(disk.SuperBlockData(bts))
	if err != nil {
		logrus.Errorf("disk.SuperBlock.Decode err:%v,", err)
		return err
	}

	// 解析 inode bitmap
	err = disk.InodeBitMap.Decode(disk.InodeBitMapData(bts), int64(disk.SuperBlock.InodeNum))
	if err != nil {
		logrus.Errorf("disk.InodeBitMap.Decode err:%v,", err)
		return err
	}

	// 解析dataBlockBitMap
	err = disk.DataBlockBitMap.Decode(disk.DataBlockBitMapData(bts), int64(disk.SuperBlock.ZoneNum))
	if err != nil {
		logrus.Errorf("DecodeDataBlockBitMap err:%v,", err)
		return err
	}

	// 解析InodeTable
	err = disk.InodeTable.Decode(disk.InodeTableData(bts), bts, int64(disk.SuperBlock.InodeNum))
	if err != nil {
		logrus.Errorf("DecodeInodeTable err:%v,", err)
		return err
	}

	return err
}
