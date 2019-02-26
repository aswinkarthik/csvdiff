package digest_test

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
)

func TestEngine_GenerateFileDigest(t *testing.T) {
	firstLine := "1,first-line,some-columne,friday"
	firstKey := xxhash.Sum64String("1")
	firstDigest := xxhash.Sum64String(firstLine)
	fridayDigest := xxhash.Sum64String("friday")

	secondLine := "2,second-line,nobody-needs-this,saturday"
	secondKey := xxhash.Sum64String("2")
	secondDigest := xxhash.Sum64String(secondLine)
	saturdayDigest := xxhash.Sum64String("saturday")

	t.Run("should create digest for given key and all values", func(t *testing.T) {
		conf := digest.Config{
			Reader: strings.NewReader(firstLine + "\n" + secondLine),
			Key:    []int{0},
		}

		engine := digest.NewEngine(conf)
		defer engine.Close()

		fd, err := engine.GenerateFileDigest()

		assert.NoError(t, err)

		expectedDigest := map[uint64]uint64{firstKey: firstDigest, secondKey: secondDigest}

		assert.Equal(t, expectedDigest, fd.Digests)
	})

	t.Run("should create digest skeeping source", func(t *testing.T) {
		conf := digest.Config{
			Reader:     strings.NewReader(firstLine + "\n" + secondLine),
			Key:        []int{0},
			KeepSource: true,
		}

		engine := digest.NewEngine(conf)
		defer engine.Close()

		fd, err := engine.GenerateFileDigest()

		assert.NoError(t, err)

		expectedDigest := map[uint64]uint64{firstKey: firstDigest, secondKey: secondDigest}
		expectedSourceMap := map[uint64][]string{
			firstKey:  strings.Split(firstLine, ","),
			secondKey: strings.Split(secondLine, ","),
		}

		assert.Equal(t, expectedDigest, fd.Digests)
		assert.Equal(t, expectedSourceMap, fd.SourceMap)
	})

	t.Run("should create digest for given key and given values", func(t *testing.T) {
		conf := digest.Config{
			Reader: strings.NewReader(firstLine + "\n" + secondLine),
			Key:    []int{0},
			Value:  []int{3},
		}

		engine := digest.NewEngine(conf)
		defer engine.Close()

		fd, err := engine.GenerateFileDigest()

		expectedDigest := map[uint64]uint64{firstKey: fridayDigest, secondKey: saturdayDigest}

		assert.NoError(t, err)
		assert.Equal(t, expectedDigest, fd.Digests)
	})

	t.Run("should return ParseError if csv reading fails", func(t *testing.T) {
		conf := digest.Config{
			Reader: strings.NewReader(firstLine + "\n" + "some-random-line"),
			Key:    []int{0},
			Value:  []int{3},
		}

		engine := digest.NewEngine(conf)
		defer engine.Close()

		fd, err := engine.GenerateFileDigest()

		assert.Error(t, err)

		_, isParseError := err.(*csv.ParseError)

		assert.True(t, isParseError)
		assert.Nil(t, fd)
	})
}
