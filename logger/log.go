package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"strconv"
)

type FileMgr struct {
	path     string
	filename string
}

func NewFileMgr(path string, filename string) *FileMgr {
	return &FileMgr{path: path, filename: filename}
}

// isExist 判断文件或文件夹是否存在
func (mgr *FileMgr) isExist() bool {
	_, err := os.Stat(mgr.path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

// CheckAndCreateDir 校验并创建文件夹
func (mgr *FileMgr) CheckAndCreateDir() error {
	//文件夹不存在
	if !mgr.isExist() {
		//递归创建文件夹
		err := os.MkdirAll(mgr.path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}

// CheckAndCreateFIle 校验并创建文件
func (mgr *FileMgr) CheckAndCreateFIle() error {
	//文件不存在
	if !mgr.isExist() {
		//递归创建文件夹
		_, err := os.Create(mgr.filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitLogger() {
	//初始化日志句柄
	logPath := "/var/log"
	logfile := logPath + "/minix_fuse.log"
	err := NewFileMgr(logPath, logfile).CheckAndCreateDir()
	if err != nil {
		panic(fmt.Sprintf("CheckAndCreateDir /var/log error. err:%v", err))
	}
	var file *os.File
	file, err = os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(fmt.Sprintf("os.OpenFile /var/log/information_flow.log error. err:%v", err))
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			//处理文件名
			fileName := path.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
			return path.Base(frame.Function), fileName
		},
	})
	logrus.SetOutput(file)
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
}
