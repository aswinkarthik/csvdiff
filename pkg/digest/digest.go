package digest

import (
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
	keyCsv := make([]string, len(keyPositions))
	for i, pos := range keyPositions {
		keyCsv[i] = csv[pos]
	}

	key := xxhash.Sum64String(strings.Join(keyCsv, ","))
	digest := xxhash.Sum64String(strings.Join(csv, ","))

	return Digest{Key: key, Value: digest}

}

type DigestConfig struct {
	KeyPositions []int
	Encoder      encoder.Encoder
	Reader       io.Reader
	Writer       io.Writer
}

func Create(config DigestConfig) (map[uint64]uint64, error) {
	reader := csv.NewReader(config.Reader)

	output := make(map[uint64]uint64)
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		digest := CreateDigest(line, config.KeyPositions)
		output[digest.Key] = digest.Value
	}

	// config.Encoder.Encode(output, config.Writer)
	return output, nil
}
