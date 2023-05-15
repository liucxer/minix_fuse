package minix_decoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
)

// InodeBitMap  3K大小
type InodeBitMap struct {
	InodeMap []bool
}

// DiskInodeBitMap inodeBitMap的磁盘结构 3K大小
type DiskInodeBitMap struct {
}

func (inodeBitMap *InodeBitMap) String() string {
	var (
		res string
	)

	res += "inodeBitMap: "
	for i := 0; i < len(inodeBitMap.InodeMap); i++ {
		if !inodeBitMap.InodeMap[i] {
			continue
		}

		res += strconv.Itoa(i+1) + ","
	}

	return res
}

func (inodeBitMap *InodeBitMap) Decode(bts []byte, inodeNum int64) error {
	buf := bytes.NewBuffer(bts)
	inodes := make([]uint8, inodeNum/8)
	err := binary.Read(buf, binary.LittleEndian, inodes)
	if err != nil {
		logrus.Errorf("binary.Read err:%v", err)
		return err
	}
	for _, item := range inodes {
		tmp := make([]bool, 8)
		if item&0x01 == 0x01 {
			tmp[0] = true
		} else {
			tmp[0] = false
		}

		if item&0x02 == 0x02 { // 0000 0010
			tmp[1] = true
		} else {
			tmp[1] = false
		}

		if item&0x04 == 0x04 { // 0000 0100
			tmp[2] = true
		} else {
			tmp[2] = false
		}

		if item&0x08 == 0x08 { // 0000 1000
			tmp[3] = true
		} else {
			tmp[3] = false
		}
		if item&0x10 == 0x10 { // 0001 0000
			tmp[4] = true
		} else {
			tmp[4] = false
		}
		if item&0x20 == 0x20 { // 0010 0000
			tmp[5] = true
		} else {
			tmp[5] = false
		}
		if item&0x40 == 0x40 { // 0100 0000
			tmp[6] = true
		} else {
			tmp[6] = false
		}
		if item&0x80 == 0x80 { // 1000 0000
			tmp[7] = true
		} else {
			tmp[7] = false
		}
		// 0000 0111
		inodeBitMap.InodeMap = append(inodeBitMap.InodeMap, tmp...)
	}

	return err
}

// Load 从磁盘加载到内存
func (inodeBitMap *InodeBitMap) Load() (int64, error) {
	return 0, nil
}

// Save 将内存数据刷到磁盘
func (inodeBitMap *InodeBitMap) Save() (int64, error) {
	return 0, nil
}

// Apply 申请inode bit
func (inodeBitMap *InodeBitMap) Apply() (int64, error) {
	for i := 0; i < len(inodeBitMap.InodeMap); i++ {
		if !inodeBitMap.InodeMap[i] {
			return int64(i), nil
		}
	}

	return 0, errors.New("no more inode bit")
}

// Release 释放inode bit
func (inodeBitMap *InodeBitMap) Release(offset int64) error {
	if len(inodeBitMap.InodeMap) < int(offset+1) {
		return fmt.Errorf("inodeBitMap Release error. length not enougth. %d < %d", len(inodeBitMap.InodeMap), int(offset+1))
	}
	inodeBitMap.InodeMap[offset] = false

	return nil
}
