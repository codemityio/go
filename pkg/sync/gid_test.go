package sync

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGIDReturnsNonZero(t *testing.T) {
	require.NotZero(t, GID(), "expected non-zero GID")
}

func TestGIDUniqueAcrossGoroutines(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)

	var gid1, gid2 int64

	go func() {
		defer wg.Done()

		gid1 = GID()
	}()

	go func() {
		defer wg.Done()

		gid2 = GID()
	}()

	wg.Wait()

	require.NotEqual(t, gid1, gid2)
}
