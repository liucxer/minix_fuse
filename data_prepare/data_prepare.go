package data_prepare

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

func CreateFile(path string, data []byte) error {
	var err error

	// 在inodeBitMap中 找到一个未使用的位置, 修改使用标志位, 标记使用
	// 在对应的inode表填写数据
	// 估算byte大小，需要占用N个块， 则申请N个标志,
	// 在对应的block中填写数据
	// 修改 inode表的 zone区域

	return err
}

func PrepareData(size int64) []byte {
	rand.Seed(time.Now().UnixNano()) // 设置随机种子

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // 字符集
	// 字符串长度
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return b
}

func BatchCreateFile(fileNum int64, filePrefix string, fileSize int64) error {
	var err error

	for i := int64(0); i < fileNum; i++ {
		fileName := filePrefix + strconv.Itoa(int(i))
		data := PrepareData(fileSize)
		err = CreateFile(fileName, data)
		if err != nil {
			logrus.Errorf("CreateFile err:%v", err)
			return err
		}
	}

	return err
}
