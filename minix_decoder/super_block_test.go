package minix_decoder_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/liucxer/minix_fuse/minix_decoder"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const DevicePath = "/dev/sdb"

func TestSuperBlock_Load(t *testing.T) {
	fd, err := os.OpenFile(DevicePath, os.O_RDWR, 0666)
	require.NoError(t, err)

	defer func() { _ = fd.Close() }()

	var superBlock minix_decoder.SuperBlock
	err = superBlock.Load(fd)
	require.NoError(t, err)
	spew.Dump(superBlock)

	superBlock.State = 1
	err = superBlock.Save(fd)
	require.NoError(t, err)
}
