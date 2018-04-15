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
	Row   string
}

// CreateDigest creates a Digest for each line of csv.
// There will be one Digest per line
func CreateDigest(csv []string, keyPositions []int) Digest {
	keyCsv := make([]string, len(keyPositions))
	for i, pos := range keyPositions {
		keyCsv[i] = csv[pos]
	}

	row := strings.Join(csv, ",")
	key := xxhash.Sum64String(strings.Join(keyCsv, ","))
	digest := xxhash.Sum64String(row)

	return Digest{Key: key, Value: digest, Row: row}

}

type DigestConfig struct {
	KeyPositions []int
	Encoder      encoder.Encoder
	Reader       io.Reader
	Writer       io.Writer
	SourceMap    bool
}

func Create(config DigestConfig) (map[uint64]uint64, map[uint64]string, error) {
	reader := csv.NewReader(config.Reader)

	output := make(map[uint64]uint64)
	sourceMap := make(map[uint64]string)
	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
		digest := CreateDigest(line, config.KeyPositions)
		output[digest.Key] = digest.Value
		if config.SourceMap {
			sourceMap[digest.Key] = digest.Row
		}
	}

	// config.Encoder.Encode(output, config.Writer)
	return output, sourceMap, nil
}
