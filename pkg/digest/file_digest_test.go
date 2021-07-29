package digest_test

import (
	"sync"
	"testing"

	"github.com/aswinkarthik/csvdiff/pkg/digest"

	"github.com/stretchr/testify/assert"
)

func TestNewFileDigest(t *testing.T) {
	fd := digest.NewFileDigest()

	assert.NotNil(t, fd)
	assert.Zero(t, len(fd.Digests))
	assert.Zero(t, len(fd.SourceMap))
}

func TestFileDigest_Append(t *testing.T) {
	fd := digest.NewFileDigest()

	fd.Append(digest.Digest{Key: uint64(1), Value: uint64(1)})

	assert.NotNil(t, fd)
	assert.Len(t, fd.Digests, 1)
	assert.Len(t, fd.SourceMap, 1)
	assert.Len(t, fd.SourceMap[uint64(1)], 1)
	assert.Len(t, fd.SourceMap[uint64(1)][0], 0)
}

func TestFileDigest_SafeAppend(t *testing.T) {
	fd := digest.NewFileDigest()

	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i uint64) {
			fd.SafeAppend(digest.Digest{Key: i, Value: i})
			wg.Done()
		}(uint64(i))
	}

	wg.Wait()
	assert.NotNil(t, fd)
	assert.Len(t, fd.Digests, 1000)
	assert.Len(t, fd.SourceMap, 1000)
}
