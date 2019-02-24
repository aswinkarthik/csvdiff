package digest

import (
	"encoding/csv"
	"fmt"
	"runtime"
	"sync"
)

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
func Diff(baseConfig, deltaConfig *Config) (Difference, error) {
	maxProcs := runtime.NumCPU()
	base, err := Create(baseConfig)

	if err != nil {
		return Difference{}, fmt.Errorf("error in base file: %v", err)
	}

	additions := make([]string, 0, len(base))
	modifications := make([]string, 0, len(base))

	messageChan := make(chan []diffMessage, bufferSize*maxProcs)
	errorChannel := make(chan error)
	defer close(errorChannel)

	go readAndCompare(base, deltaConfig, messageChan, errorChannel)

	for msgs := range messageChan {
		for _, msg := range msgs {
			if msg._type == addition {
				additions = append(additions, msg.value)
			} else if msg._type == modification {
				modifications = append(modifications, msg.value)
			}
		}
	}

	if err := <-errorChannel; err != nil {
		return Difference{}, fmt.Errorf("error in delta file: %v", err)
	}

	return Difference{Additions: additions, Modifications: modifications}, nil
}

func readAndCompare(base map[uint64]uint64, config *Config, msgChannel chan<- []diffMessage, errorChannel chan<- error) {
	reader := csv.NewReader(config.Reader)
	var wg sync.WaitGroup
	for {
		lines, eofReached, err := getNextNLines(reader)

		if err != nil {
			wg.Wait()
			close(msgChannel)
			errorChannel <- err
			return
		}

		wg.Add(1)
		go compareDigestForNLines(base, lines, config, msgChannel, &wg)

		if eofReached {
			break
		}
	}
	wg.Wait()
	close(msgChannel)
	errorChannel <- nil
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
				value := config.Include.MapToValue(line)
				// Modification
				output[diffCounter] = diffMessage{
					value: value,
					_type: modification,
				}
				diffCounter++
			}
		} else {
			value := config.Include.MapToValue(line)
			// Not present in base. So Addition.
			output[diffCounter] = diffMessage{
				value: value,
				_type: addition,
			}
			diffCounter++
		}
	}

	msgChannel <- output[:diffCounter]
	wg.Done()
}
