package digest_test

import (
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

	expectedDigest := digest.Digest{Key: firstKey, Value: firstLineDigest}

	actualDigest := digest.CreateDigest(strings.Split(firstLine, digest.Separator), []int{0}, []int{})

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

	testConfig := &digest.Config{
		Reader: strings.NewReader(firstLine + "\n" + secondLine),
		Key:    []int{0},
	}

	actualDigest := digest.Create(testConfig)

	expectedDigest := map[uint64]uint64{firstKey: firstDigest, secondKey: secondDigest}

	assert.Equal(t, expectedDigest, actualDigest)

	testConfig = &digest.Config{
		Reader: strings.NewReader(firstLine + "\n" + secondLine),
		Key:    []int{0},
		Value:  []int{3},
	}

	actualDigest = digest.Create(testConfig)
	expectedDigest = map[uint64]uint64{firstKey: fridayDigest, secondKey: saturdayDigest}

	assert.Equal(t, expectedDigest, actualDigest)
}
