package digest

import (
	"encoding/csv"
	"io"
	"strings"

	"github.com/cespare/xxhash"
)

const Separator = ","

// Digest represents the binding of the key of each csv line
// and the digest that gets created for the entire line
type Digest struct {
	Key   uint64
	Value uint64
	Row   string
}

// CreateDigest creates a Digest for each line of csv.
// There will be one Digest per line
func CreateDigest(csv []string, pKey Positions, pRow Positions) Digest {
	row := strings.Join(csv, Separator)
	key := xxhash.Sum64String(pKey.MapToValue(csv))
	digest := xxhash.Sum64String(pRow.MapToValue(csv))

	return Digest{Key: key, Value: digest, Row: row}

}

type Config struct {
	KeyPositions []int
	Key          Positions
	Value        Positions
	Reader       io.Reader
	Writer       io.Writer
	SourceMap    bool
}

func NewConfig(r io.Reader, createSourceMap bool, primaryKey Positions, valueColumns Positions) *Config {
	return &Config{
		Reader:    r,
		SourceMap: createSourceMap,
		Key:       primaryKey,
		Value:     valueColumns,
	}
}

func Create(config *Config) (map[uint64]uint64, map[uint64]string, error) {
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
		digest := CreateDigest(line, config.Key, config.Value)
		output[digest.Key] = digest.Value
		if config.SourceMap {
			sourceMap[digest.Key] = digest.Row
		}
	}

	return output, sourceMap, nil
}
