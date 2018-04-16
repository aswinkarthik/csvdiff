package digest

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
)

func TestCreateDigest(t *testing.T) {
	firstLine := "1,someline"
	firstKey := xxhash.Sum64String("1")
	firstLineDigest := xxhash.Sum64String(firstLine)

	expectedDigest := Digest{Key: firstKey, Value: firstLineDigest, Row: firstLine}

	actualDigest := CreateDigest(strings.Split(firstLine, Separator), []int{0}, []int{})

	assert.Equal(t, expectedDigest, actualDigest)
}

func TestDigestForFile(t *testing.T) {
	firstLine := "1,first-line,some-columne,friday"
	firstKey := xxhash.Sum64String("1")
	firstDigest := xxhash.Sum64String(firstLine)

	secondLine := "2,second-line,nobody-needs-this,saturday"
	secondKey := xxhash.Sum64String("2")
	secondDigest := xxhash.Sum64String(secondLine)

	var outputBuffer bytes.Buffer

	testConfig := &Config{
		Reader:       strings.NewReader(firstLine + "\n" + secondLine),
		Writer:       &outputBuffer,
		KeyPositions: []int{0},
		Key:          []int{0},
		SourceMap:    true,
	}

	actualDigest, sourceMap, err := Create(testConfig)

	expectedDigest := map[uint64]uint64{firstKey: firstDigest, secondKey: secondDigest}
	expectedSourceMap := map[uint64]string{firstKey: firstLine, secondKey: secondLine}

	assert.Nil(t, err, "error at DigestForFile")
	assert.Equal(t, expectedDigest, actualDigest)
	assert.Equal(t, expectedSourceMap, sourceMap)

	// No source map
	testConfigWithoutSourceMap := &Config{
		Reader:       strings.NewReader(firstLine + "\n" + secondLine),
		Writer:       &outputBuffer,
		KeyPositions: []int{0},
		Key:          []int{0},
		SourceMap:    false,
	}

	actualDigest, sourceMap, err = Create(testConfigWithoutSourceMap)

	assert.Nil(t, err, "error at DigestForFile")
	assert.Equal(t, expectedDigest, actualDigest)
	assert.Equal(t, map[uint64]string{}, sourceMap)
}
