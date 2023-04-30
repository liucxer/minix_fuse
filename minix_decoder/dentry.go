package minix_decoder

import "fmt"

/*
// 文件目录项结构。
struct dir_entry
{
  unsigned short inode;		// i 节点。
  char name[NAME_LEN];		// 文件名。
};
*/

const NameLen = 30

type DirEntry struct {
	InodeNo uint16
	Name    [NameLen]byte
}

func (m DirEntry) String() string {
	var bts []byte
	for _, b := range m.Name {
		if b == 0 {
			break
		}
		bts = append(bts, byte(b))
	}

	return fmt.Sprintf("DirEntry inodeNo:%d, Name:%s", m.InodeNo, string(bts))
}
