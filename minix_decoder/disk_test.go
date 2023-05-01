package minix_decoder_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/liucxer/minix_fuse/minix_decoder"
	"github.com/stretchr/testify/require"
	"testing"
)

type File struct {
	Inode minix_decoder.Inode
	Path  string
	Data  string
}

func TestDiskData_Decode(t *testing.T) {
	files, err := minix_decoder.GetFiles("/dev/vdb")
	spew.Dump(files)
	require.NoError(t, err)
}
