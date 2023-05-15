package data_prepare_test

import (
	"github.com/liucxer/minix_fuse/data_prepare"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBatchCreateFile(t *testing.T) {
	for i := 0; i < 100; i++ {
		err := data_prepare.BatchCreateFile(10, "AAA", int64(i)*1024)
		require.NoError(t, err)
	}
}
