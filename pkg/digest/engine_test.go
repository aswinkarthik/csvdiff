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

		dChan, eChan := engine.StreamDigests()

		err := <-eChan
		assert.NoError(t, err)

		actualDigest := digestsFrom(dChan)
		expectedDigest := []digest.Digest{
			{Key: firstKey, Value: firstDigest},
			{Key: secondKey, Value: secondDigest},
		}

		assert.ElementsMatch(t, expectedDigest, actualDigest)
	})

	t.Run("should create digest skeeping source", func(t *testing.T) {
		conf := digest.Config{
			Reader:     strings.NewReader(firstLine + "\n" + secondLine),
			Key:        []int{0},
			KeepSource: true,
		}

		engine := digest.NewEngine(conf)

		dChan, eChan := engine.StreamDigests()

		err := <-eChan
		assert.NoError(t, err)

		actualDigest := digestsFrom(dChan)
		expectedDigest := []digest.Digest{
			{Key: firstKey, Value: firstDigest, Source: strings.Split(firstLine, ",")},
			{Key: secondKey, Value: secondDigest, Source: strings.Split(secondLine, ",")},
		}

		assert.ElementsMatch(t, expectedDigest, actualDigest)
	})

	t.Run("should create digest for given key and given values", func(t *testing.T) {
		conf := digest.Config{
			Reader: strings.NewReader(firstLine + "\n" + secondLine),
			Key:    []int{0},
			Value:  []int{3},
		}

		engine := digest.NewEngine(conf)

		dChan, eChan := engine.StreamDigests()

		err := <-eChan
		assert.NoError(t, err)

		actualDigest := digestsFrom(dChan)
		expectedDigest := []digest.Digest{
			{Key: firstKey, Value: fridayDigest},
			{Key: secondKey, Value: saturdayDigest},
		}

		assert.ElementsMatch(t, expectedDigest, actualDigest)
	})

	t.Run("should return ParseError if csv reading fails", func(t *testing.T) {
		conf := digest.Config{
			Reader: strings.NewReader(firstLine + "\n" + "some-random-line"),
			Key:    []int{0},
			Value:  []int{3},
		}

		engine := digest.NewEngine(conf)

		dChan, eChan := engine.StreamDigests()

		err := <-eChan

		assert.Error(t, err)

		_, isParseError := err.(*csv.ParseError)

		assert.True(t, isParseError)

		actualDigest := digestsFrom(dChan)
		assert.Empty(t, actualDigest)
	})
}

func digestsFrom(digestChan chan []digest.Digest) []digest.Digest {
	result := make([]digest.Digest, 0, 10)

	for rcvd := range digestChan {
		for _, d := range rcvd {
			result = append(result, d)
		}
	}

	return result
}
