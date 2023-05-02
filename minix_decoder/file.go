package minix_decoder

import "github.com/sirupsen/logrus"

type File struct {
	Inode   Inode
	Path    string
	Data    string
	InodeNo int64
	Files   []File
}

func GetFileMap(devicePath string) (map[int64]File, error) {
	var (
		disk DiskData
		err  error
	)

	err = disk.Decode(devicePath)

	inodeMap := map[int64]string{}
	inodeMap[0] = "/"
	for _, inodeItem := range disk.InodeTable.InodeItems {
		if inodeItem.Mode.IsDir() {
			for _, dir := range inodeItem.DirEntry {
				_, ok := inodeMap[int64(dir.InodeNo)]
				if ok {
					if dir.String() != "." && dir.String() != ".." {
						inodeMap[int64(dir.InodeNo)] = dir.String()
					}
				} else {
					inodeMap[int64(dir.InodeNo)] = dir.String()
				}
			}
		}
	}

	var files []File
	for i := 0; i < len(disk.InodeTable.InodeItems); i++ {
		tmpInodeItem := disk.InodeTable.InodeItems[i]
		if !tmpInodeItem.Mode.IsDir() && !tmpInodeItem.Mode.IsReg() {
			continue
		}

		var file File
		file.Inode = tmpInodeItem.Inode
		file.Path = inodeMap[int64(i+1)]
		file.Data = string(tmpInodeItem.Data)
		file.InodeNo = int64(i + 1)

		var subFiles []File
		for _, dentry := range tmpInodeItem.DirEntry {
			var subFile File
			subFile.InodeNo = int64(dentry.InodeNo)
			subFile.Path = dentry.String()
			subFiles = append(subFiles, subFile)
		}
		file.Files = subFiles
		files = append(files, file)
	}

	fileMap := make(map[int64]File)
	for _, file := range files {
		fileMap[file.InodeNo] = file
	}
	logrus.Infof("fileMap:%+v", fileMap)
	return fileMap, err
}
