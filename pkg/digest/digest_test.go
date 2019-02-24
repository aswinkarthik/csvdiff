package digest_test

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
)

func TestCreateDigest(t *testing.T) {
	firstLine := "1,someline"
	firstKey := xxhash.Sum64String("1")
	firstLineDigest := xxhash.Sum64String(firstLine)

	expectedDigest := digest.Digest{Key: firstKey, Value: firstLineDigest, Source: nil}

	actualDigest := digest.CreateDigest(strings.Split(firstLine, digest.Separator), []int{0}, []int{})

	assert.Equal(t, expectedDigest, actualDigest)
}

func TestCreateDigestWithSource(t *testing.T) {
	firstLine := "1,someline"
	firstKey := xxhash.Sum64String("1")
	firstLineDigest := xxhash.Sum64String(firstLine)

	expectedDigest := digest.Digest{
		Key:    firstKey,
		Value:  firstLineDigest,
		Source: strings.Split(firstLine, ","),
	}

	actualDigest := digest.CreateDigestWithSource(strings.Split(firstLine, digest.Separator), []int{0}, []int{})

	assert.Equal(t, expectedDigest, actualDigest)
}

func TestDigestForFile(t *testing.T) {
	firstLine := "1,first-line,some-columne,friday"
	firstKey := xxhash.Sum64String("1")
	firstDigest := xxhash.Sum64String(firstLine)
	fridayDigest := xxhash.Sum64String("friday")

	secondLine := "2,second-line,nobody-needs-this,saturday"
	secondKey := xxhash.Sum64String("2")
	secondDigest := xxhash.Sum64String(secondLine)
	saturdayDigest := xxhash.Sum64String("saturday")

	t.Run("should create digest for given key and all values", func(t *testing.T) {
		testConfig := &digest.Config{
			Reader: strings.NewReader(firstLine + "\n" + secondLine),
			Key:    []int{0},
		}

		actualDigest, err := digest.Create(testConfig)

		expectedDigest := map[uint64]uint64{firstKey: firstDigest, secondKey: secondDigest}

		assert.NoError(t, err)
		assert.Equal(t, expectedDigest, actualDigest)
	})

	t.Run("should create digest for given key and given values", func(t *testing.T) {
		testConfig := &digest.Config{
			Reader: strings.NewReader(firstLine + "\n" + secondLine),
			Key:    []int{0},
			Value:  []int{3},
		}

		actualDigest, err := digest.Create(testConfig)
		expectedDigest := map[uint64]uint64{firstKey: fridayDigest, secondKey: saturdayDigest}

		assert.NoError(t, err)
		assert.Equal(t, expectedDigest, actualDigest)
	})

	t.Run("should return ParseError if csv reading fails", func(t *testing.T) {
		testConfig := &digest.Config{
			Reader: strings.NewReader(firstLine + "\n" + "some-random-line"),
			Key:    []int{0},
			Value:  []int{3},
		}

		actualDigest, err := digest.Create(testConfig)

		assert.Error(t, err)

		_, isParseError := err.(*csv.ParseError)

		assert.True(t, isParseError)
		assert.Nil(t, actualDigest)
	})
}

func TestNewConfig(t *testing.T) {
	r := strings.NewReader("a,csv,as,str")
	primaryColumns := digest.Positions{0}
	values := digest.Positions{0, 1, 2}
	include := digest.Positions{0, 1}
	keepSource := true

	t.Run("should create config from given params", func(t *testing.T) {
		conf := digest.NewConfig(r, primaryColumns, values, include, keepSource)
		expectedConf := digest.Config{
			Reader:     r,
			Key:        primaryColumns,
			Value:      values,
			Include:    include,
			KeepSource: keepSource,
		}

		assert.Equal(t, expectedConf, *conf)
	})

	t.Run("should use valueColumns as includeColumns for includes not specified", func(t *testing.T) {
		conf := digest.NewConfig(r, primaryColumns, values, nil, keepSource)
		expectedConf := digest.Config{
			Reader:     r,
			Key:        primaryColumns,
			Value:      values,
			Include:    values,
			KeepSource: keepSource,
		}

		assert.Equal(t, expectedConf, *conf)
	})
}
