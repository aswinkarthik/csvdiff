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
	Key   uint64
	Value uint64
}

// CreateDigest creates a Digest for each line of csv.
// There will be one Digest per line
func CreateDigest(csv []string, pKey Positions, pRow Positions) Digest {
	key := xxhash.Sum64String(pKey.MapToValue(csv))
	digest := xxhash.Sum64String(pRow.MapToValue(csv))

	return Digest{Key: key, Value: digest}

}

// Config represents configurations that can be passed
// to create a Digest.
type Config struct {
	Key    Positions
	Value  Positions
	Reader io.Reader
}

// NewConfig creates an instance of Config struct.
func NewConfig(r io.Reader, primaryKey Positions, valueColumns Positions) *Config {
	return &Config{
		Reader: r,
		Key:    primaryKey,
		Value:  valueColumns,
	}
}

const bufferSize = 512

// Create can create a Digest using the Configurations passed.
// It returns the digest as a map[uint64]uint64.
// It can also keep track of the Source line.
func Create(config *Config) (map[uint64]uint64, error) {
	maxProcs := runtime.NumCPU()
	reader := csv.NewReader(config.Reader)

	output := make(map[uint64]uint64)

	digestChannel := make(chan []Digest, bufferSize*maxProcs)
	errorChannel := make(chan error)
	defer close(errorChannel)

	go readAndProcess(config, reader, digestChannel, errorChannel)

	for digests := range digestChannel {
		for _, digest := range digests {
			output[digest.Key] = digest.Value
		}
	}

	if err := <-errorChannel; err != nil {
		return nil, err
	}

	return output, nil
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

	for i, line := range lines {
		output[i] = CreateDigest(line, config.Key, config.Value)
	}

	digestChannel <- output
	wg.Done()
}
