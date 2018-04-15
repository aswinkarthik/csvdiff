package digest

import (
	"bytes"
	"encoding/csv"
	"io"
	"strings"

	"github.com/aswinkarthik93/csv-digest/pkg/encoder"
	"github.com/cespare/xxhash"
)

// CsvDigest represents the binding of the key of each csv line
// and the digest that gets created for the entire line
type CsvDigest struct {
	Key    uint64
	Digest uint64
}

// CreateDigest creates a Digest for each line of csv.
// There will be one CsvDigest per line
func CreateDigest(csv []string, keyPositions []int) CsvDigest {
	var keyBuffer bytes.Buffer
	return CreateDigestWithBuffer(csv, keyPositions, &keyBuffer)
}

// CreateDigestWithBuffer creates a Digest for each line of csv.
// Also takes a buffer which can be passed to optimize on allocating a buffer for
// computing digest of the key
func CreateDigestWithBuffer(csv []string, keyPositions []int, b *bytes.Buffer) CsvDigest {
	for _, pos := range keyPositions {
		b.WriteString(csv[pos])
	}

	key := xxhash.Sum64(b.Bytes())
	digest := xxhash.Sum64String(strings.Join(csv, ","))

	b.Reset()
	return CsvDigest{Key: key, Digest: digest}

}

type DigestConfig struct {
	KeyPositions []int
	Encoder      encoder.Encoder
	Reader       io.Reader
	Writer       io.Writer
}

func DigestForFile(config DigestConfig) error {
	reader := csv.NewReader(config.Reader)
	for {
		line, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		config.Encoder.Encode(CreateDigest(line, config.KeyPositions), config.Writer)
	}

	return nil
}
