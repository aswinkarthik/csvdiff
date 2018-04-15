package digest

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"github.com/aswinkarthik93/csv-digest/pkg/encoder"
	"github.com/cespare/xxhash"
)

// Digest represents the binding of the key of each csv line
// and the digest that gets created for the entire line
type Digest struct {
	Key   uint64
	Value uint64
}

// CreateDigest creates a Digest for each line of csv.
// There will be one Digest per line
func CreateDigest(csv []string, keyPositions []int) Digest {
	var keyBuffer bytes.Buffer
	return CreateDigestWithBuffer(csv, keyPositions, &keyBuffer)
}

// CreateDigestWithBuffer creates a Digest for each line of csv.
// Also takes a buffer which can be passed to optimize on allocating a buffer for
// computing digest of the key
func CreateDigestWithBuffer(csv []string, keyPositions []int, b *bytes.Buffer) Digest {
	for _, pos := range keyPositions {
		b.WriteString(csv[pos])
	}

	key := xxhash.Sum64(b.Bytes())
	digest := xxhash.Sum64String(strings.Join(csv, ","))

	b.Reset()
	return Digest{Key: key, Value: digest}

}

type DigestConfig struct {
	KeyPositions []int
	Encoder      encoder.Encoder
	Reader       io.Reader
	Writer       io.Writer
}

func DigestForFile(config DigestConfig) error {
	bufferedReader := bufio.NewReader(config.Reader)
	reader := csv.NewReader(bufferedReader)
	lines, err := reader.ReadAll()
	output := make([]Digest, len(lines))
	for i, line := range lines {

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		output[i] = CreateDigest(line, config.KeyPositions)
	}

	config.Encoder.Encode(toHash(output), config.Writer)
	return nil
}

func toHash(digests []Digest) map[uint64]uint64 {
	result := make(map[uint64]uint64, len(digests))

	for _, digest := range digests {
		result[digest.Key] = digest.Value
	}

	return result
}
