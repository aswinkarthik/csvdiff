package digest

import (
	"fmt"
	"runtime"
)

type messageType int

const (
	addition     messageType = iota
	modification messageType = iota
)

// Differences represents the differences
// between 2 csv content
type Differences struct {
	Additions     []Addition
	Modifications []Modification
}

// Addition is a row appearing in delta but missing in base
type Addition []string

// Modification is a row present in both delta and base
// with the values column changed in delta
type Modification struct {
	Original []string
	Current  []string
}

type message struct {
	original []string
	current  []string
	_type    messageType
}

// Diff finds the Differences between baseConfig and deltaConfig
func Diff(baseConfig, deltaConfig Config) (Differences, error) {
	baseEngine := NewEngine(baseConfig)
	baseDigestChannel, baseErrorChannel := baseEngine.StreamDigests()

	baseFileDigest := NewFileDigest()
	for digests := range baseDigestChannel {
		for _, d := range digests {
			baseFileDigest.Append(d)
		}
	}

	if err := <-baseErrorChannel; err != nil {
		return Differences{}, fmt.Errorf("error processing base file: %v", err)
	}

	deltaConfig.KeepSource = true
	deltaEngine := NewEngine(deltaConfig)
	deltaDigestChannel, deltaErrorChannel := deltaEngine.StreamDigests()

	additions := make([]Addition, 0)
	modifications := make([]Modification, 0)

	msgChannel := streamDifferences(baseFileDigest, deltaDigestChannel)
	for msg := range msgChannel {
		switch msg._type {
		case addition:
			additions = append(additions, msg.current)
		case modification:
			modifications = append(modifications, Modification{Original: msg.original, Current: msg.current})
		}
	}

	if err := <-deltaErrorChannel; err != nil {
		return Differences{}, fmt.Errorf("error processing delta file: %v", err)
	}

	return Differences{Additions: additions, Modifications: modifications}, nil
}

func streamDifferences(baseFileDigest *FileDigest, digestChannel chan []Digest) chan message {
	maxProcs := runtime.NumCPU()
	msgChannel := make(chan message, maxProcs*bufferSize)

	go func(base *FileDigest, digestChannel chan []Digest, msgChannel chan message) {
		defer close(msgChannel)

		for digests := range digestChannel {
			for _, d := range digests {
				if baseValue, present := base.Digests[d.Key]; present {
					if baseValue != d.Value {
						// Modification
						msgChannel <- message{_type: modification, current: d.Source, original: base.SourceMap[d.Key]}
					}
				} else {
					// Addition
					msgChannel <- message{_type: addition, current: d.Source}
				}
			}
		}

	}(baseFileDigest, digestChannel, msgChannel)

	return msgChannel
}
