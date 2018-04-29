package digest

import (
	"encoding/csv"
	"runtime"
	"strings"
	"sync"
)

// Compare compares two Digest maps and returns the additions and modification
// keys as arrays.
func Compare(baseDigest, newDigest map[uint64]uint64) (additions []uint64, modifications []uint64) {
	maxSize := len(newDigest)
	additions = make([]uint64, maxSize)
	modifications = make([]uint64, maxSize)

	additionCounter := 0
	modificationCounter := 0
	for k, newVal := range newDigest {
		if oldVal, present := baseDigest[k]; present {
			if newVal != oldVal {
				//Modifications
				modifications[modificationCounter] = k
				modificationCounter++
			}
		} else {
			//Additions
			additions[additionCounter] = k
			additionCounter++
		}
	}
	return additions[:additionCounter], modifications[:modificationCounter]
}

// Difference represents the additions and modifications
// between the two Configs
type Difference struct {
	Additions     []string
	Modifications []string
}

type messageType int

const (
	addition     messageType = iota
	modification messageType = iota
)

type diffMessage struct {
	_type messageType
	value string
}

// Diff will differentiate between two given configs
func Diff(baseConfig, deltaConfig *Config) Difference {
	maxProcs := runtime.NumCPU()
	base := Create(baseConfig)

	additions := make([]string, 0, len(base))
	modifications := make([]string, 0, len(base))

	messageChan := make(chan []diffMessage, bufferSize*maxProcs)

	go readAndCompare(base, deltaConfig, messageChan)

	for msgs := range messageChan {
		for _, msg := range msgs {
			if msg._type == addition {
				additions = append(additions, msg.value)
			} else if msg._type == modification {
				modifications = append(modifications, msg.value)
			}
		}
	}

	return Difference{Additions: additions, Modifications: modifications}
}

func readAndCompare(base map[uint64]uint64, config *Config, msgChannel chan<- []diffMessage) {
	reader := csv.NewReader(config.Reader)
	var wg sync.WaitGroup
	for {
		lines, eofReached := getNextNLines(reader)
		wg.Add(1)
		go compareDigestForNLines(base, lines, config, msgChannel, &wg)

		if eofReached {
			break
		}
	}
	wg.Wait()
	close(msgChannel)
}

func compareDigestForNLines(base map[uint64]uint64,
	lines [][]string,
	config *Config,
	msgChannel chan<- []diffMessage,
	wg *sync.WaitGroup,
) {
	output := make([]diffMessage, len(lines))
	diffCounter := 0
	for _, line := range lines {
		digest := CreateDigest(line, config.Key, config.Value)
		if baseValue, present := base[digest.Key]; present {
			// Present in both base and delta
			if baseValue != digest.Value {
				// Modification
				output[diffCounter] = diffMessage{
					value: strings.Join(line, Separator),
					_type: modification,
				}
				diffCounter++
			}
		} else {
			// Not present in base. So Addition.
			output[diffCounter] = diffMessage{
				value: strings.Join(line, Separator),
				_type: addition,
			}
			diffCounter++
		}
	}

	msgChannel <- output[:diffCounter]
	wg.Done()
}
