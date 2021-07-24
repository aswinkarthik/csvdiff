package digest

import (
	"fmt"
	"runtime"
)

type messageType int

const (
	addition     messageType = iota
	modification messageType = iota
	deletion     messageType = iota
)

// Differences represents the differences
// between 2 csv content
type Differences struct {
	Additions     []Addition
	Modifications []Modification
	Deletions     []Deletion
}

// Addition is a row appearing in delta but missing in base
type Addition []string

// Deletion is a row appearing in base but missing in delta
type Deletion []string

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

	deltaEngine := NewEngine(deltaConfig)
	deltaDigestChannel, deltaErrorChannel := deltaEngine.StreamDigests()

	additions := make([]Addition, 0)
	modifications := make([]Modification, 0)
	deletions := make([]Deletion, 0)

	msgChannel := streamDifferences(baseFileDigest, deltaDigestChannel)
	for msg := range msgChannel {
		switch msg._type {
		case addition:
			additions = append(additions, msg.current)
		case modification:
			modifications = append(modifications, Modification{Original: msg.original, Current: msg.current})
		case deletion:
			deletions = append(deletions, msg.current)
		default:
			continue
		}
	}

	if err := <-deltaErrorChannel; err != nil {
		return Differences{}, fmt.Errorf("error processing delta file: %v", err)
	}

	return Differences{Additions: additions, Modifications: modifications, Deletions: deletions}, nil
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
						msgChannel <- message{_type: modification, current: d.Source, original: base.SourceMap[d.Key][0]}
					}
					// delete from sourceMap so that at the end only deletions are left in base
					sources := base.SourceMap[d.Key]
					if len(sources) == 1 {
						delete(base.SourceMap, d.Key)
						continue
					}
					sources = sources[:len(sources)-1]
					base.SourceMap[d.Key] = sources
				} else {
					// Addition
					msgChannel <- message{_type: addition, current: d.Source}
				}
			}
		}
		for _, value := range base.SourceMap {
			msgChannel <- message{_type: deletion, current: value[0]}
		}

	}(baseFileDigest, digestChannel, msgChannel)

	return msgChannel
}
