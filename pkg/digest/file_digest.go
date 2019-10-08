package digest

import (
	"sync"
)

// FileDigest represents the digests created from one file
type FileDigest struct {
	Digests   map[uint64]uint64
	SourceMap map[uint64][]string
	lock      *sync.Mutex
}

// NewFileDigest to instantiate a new FileDigest
func NewFileDigest() *FileDigest {
	return &FileDigest{
		Digests:   make(map[uint64]uint64),
		SourceMap: make(map[uint64][]string),
		lock:      &sync.Mutex{},
	}
}

// Append a Digest to a FileDigest
// This operation is not thread safe
func (f *FileDigest) Append(d Digest) {
	f.Digests[d.Key] = d.Value
	f.SourceMap[d.Key] = d.Source
}

// SafeAppend a Digest to a FileDigest
// This operation is thread safe
func (f *FileDigest) SafeAppend(d Digest) {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.Digests[d.Key] = d.Value
	f.SourceMap[d.Key] = d.Source
}
