package digest

import (
	"encoding/csv"
	"io"
	"runtime"
	"sync"

	"github.com/cespare/xxhash"
)

// Separator for CSV. Not configurable for now.
const Separator = ","

// Digest represents the binding of the key of each csv line
// and the digest that gets created for the entire line
type Digest struct {
	Key    uint64
	Value  uint64
	Source []string
}

// CreateDigest creates a Digest for each line of csv.
// There will be one Digest per line
func CreateDigest(csv []string, pKey Positions, pRow Positions) Digest {
	key := xxhash.Sum64String(pKey.MapToValue(csv))
	digest := xxhash.Sum64String(pRow.MapToValue(csv))

	return Digest{Key: key, Value: digest}
}

// CreateDigestWithSource creates a Digest for each line of csv.
// There will be one Digest per line
func CreateDigestWithSource(csv []string, pKey Positions, pRow Positions) Digest {
	key := xxhash.Sum64String(pKey.MapToValue(csv))
	digest := xxhash.Sum64String(pRow.MapToValue(csv))

	return Digest{Key: key, Value: digest, Source: csv}
}

// Config represents configurations that can be passed
// to create a Digest.
//
// Key: The primary key positions
// Value: The Value positions that needs to be compared for diff
// Include: Include these positions in output. It is Value positions by default.
// KeepSource: return the source and target string if diff is computed
type Config struct {
	Key        Positions
	Value      Positions
	Include    Positions
	Reader     io.Reader
	KeepSource bool
}

// NewConfig creates an instance of Config struct.
func NewConfig(
	r io.Reader,
	primaryKey Positions,
	valueColumns Positions,
	includeColumns Positions,
	keepSource bool,
) *Config {
	if len(includeColumns) == 0 {
		includeColumns = valueColumns
	}

	return &Config{
		Reader:     r,
		Key:        primaryKey,
		Value:      valueColumns,
		Include:    includeColumns,
		KeepSource: keepSource,
	}
}

const bufferSize = 512

// Create can create a Digest using the Configurations passed.
// It returns the digest as a map[uint64]uint64.
// It can also keep track of the Source line.
func Create(config *Config) (map[uint64]uint64, map[uint64][]string, error) {
	maxProcs := runtime.NumCPU()
	reader := csv.NewReader(config.Reader)

	output := make(map[uint64]uint64)

	var sourceMap map[uint64][]string

	if config.KeepSource {
		sourceMap = make(map[uint64][]string)
	}

	digestChannel := make(chan []Digest, bufferSize*maxProcs)
	errorChannel := make(chan error)
	defer close(errorChannel)

	go readAndProcess(config, reader, digestChannel, errorChannel)

	for digests := range digestChannel {
		for _, digest := range digests {
			output[digest.Key] = digest.Value

			if config.KeepSource {
				sourceMap[digest.Key] = digest.Source
			}
		}
	}

	if err := <-errorChannel; err != nil {
		return nil, nil, err
	}

	return output, sourceMap, nil
}

func readAndProcess(config *Config, reader *csv.Reader, digestChannel chan<- []Digest, errorChannel chan<- error) {
	var wg sync.WaitGroup
	for {
		lines, eofReached, err := getNextNLines(reader)
		if err != nil {
			wg.Wait()
			close(digestChannel)
			errorChannel <- err
			return
		}

		wg.Add(1)
		go createDigestForNLines(lines, config, digestChannel, &wg)

		if eofReached {
			break
		}
	}
	wg.Wait()
	close(digestChannel)
	errorChannel <- nil
}

func createDigestForNLines(lines [][]string,
	config *Config,
	digestChannel chan<- []Digest,
	wg *sync.WaitGroup,
) {
	output := make([]Digest, len(lines))
	var createDigestFunc func(csv []string, pKey Positions, pRow Positions) Digest

	if config.KeepSource {
		createDigestFunc = CreateDigestWithSource
	} else {
		createDigestFunc = CreateDigest
	}

	for i, line := range lines {
		output[i] = createDigestFunc(line, config.Key, config.Value)
	}

	digestChannel <- output
	wg.Done()
}
