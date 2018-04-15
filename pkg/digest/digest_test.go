package digest

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/aswinkarthik93/csv-digest/pkg/encoder"
	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
)

func TestCreateDigest(t *testing.T) {
	firstLine := "1,someline"
	firstKey := xxhash.Sum64String("1")
	firstLineDigest := xxhash.Sum64String(firstLine)

	expectedDigest := Digest{Key: firstKey, Value: firstLineDigest}

	actualDigest := CreateDigest(strings.Split(firstLine, ","), []int{0})

	assert.Equal(t, expectedDigest, actualDigest)
}

func TestDigestForFile(t *testing.T) {
	firstLine := "1,first-line"
	firstKey := xxhash.Sum64String("1")
	firstDigest := xxhash.Sum64String(firstLine)

	secondLine := "2,second-line"
	secondKey := xxhash.Sum64String("2")
	secondDigest := xxhash.Sum64String(secondLine)

	var outputBuffer bytes.Buffer

	testConfig := DigestConfig{
		Reader:       strings.NewReader(firstLine + "\n" + secondLine),
		Writer:       &outputBuffer,
		Encoder:      encoder.JsonEncoder{},
		KeyPositions: []int{0},
	}

	err := DigestForFile(testConfig)

	actualDigest := outputBuffer.String()
	expectedDigest := fmt.Sprintf(`{"%d":%d,"%d":%d}`, firstKey, firstDigest, secondKey, secondDigest)

	assert.Nil(t, err, "error at DigestForFile")
	assert.Equal(t, expectedDigest, actualDigest)
}
