package digest

import (
	"os"
	"strings"
	"testing"

	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
)

func TestCreateDigest(t *testing.T) {
	firstLine := "1,someline"
	firstKey := xxhash.Sum64String("1")
	firstLineDigest := xxhash.Sum64String(firstLine)

	expectedDigest := Digest{Key: firstKey, Value: firstLineDigest}

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

	testConfig := &Config{
		Reader: strings.NewReader(firstLine + "\n" + secondLine),
		Key:    []int{0},
	}

	actualDigest, _, err := Create(testConfig)

	expectedDigest := map[uint64]uint64{firstKey: firstDigest, secondKey: secondDigest}

	assert.Nil(t, err, "error at DigestForFile")
	assert.Equal(t, expectedDigest, actualDigest)
}

func TestCreatePerformance(t *testing.T) {
	file, err := os.Open("../../benchmark/majestic_million.csv")
	defer file.Close()
	assert.NoError(t, err)

	config := &Config{
		Reader: file,
		Key:    []int{},
	}

	result, _, _ := Create(config)

	assert.Equal(t, 998390, len(result))
}
