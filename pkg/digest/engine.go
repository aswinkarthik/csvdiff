package digest

import (
	"encoding/csv"
	"runtime"
	"sync"
)

// Engine to create a FileDigest
type Engine struct {
	config        Config
	reader        *csv.Reader
	lock          *sync.Mutex
	digestChannel chan []Digest
	errorChannel  chan error
}

// NewEngine instantiates an engine
func NewEngine(config Config) *Engine {
	maxProcs := runtime.NumCPU()
	digestChannel := make(chan []Digest, bufferSize*maxProcs)
	errorChannel := make(chan error)

	return &Engine{
		config:        config,
		lock:          &sync.Mutex{},
		digestChannel: digestChannel,
		errorChannel:  errorChannel,
	}
}

// Close the engine after use
func (e Engine) Close() {
	close(e.errorChannel)
}

// GenerateFileDigest generates FileDigest with thread safety
func (e Engine) GenerateFileDigest() (*FileDigest, error) {
	e.lock.Lock()
	defer e.lock.Unlock()

	fd := NewFileDigest()

	var appendFunc func(Digest)

	if e.config.KeepSource {
		appendFunc = fd.Append
	} else {
		appendFunc = fd.AppendWithoutSource
	}

	go e.createDigestsInBackground()

	for digests := range e.digestChannel {
		for _, digest := range digests {
			appendFunc(digest)
		}
	}

	if err := <-e.errorChannel; err != nil {
		return nil, err
	}

	return fd, nil
}

func (e Engine) createDigestsInBackground() {
	wg := &sync.WaitGroup{}
	reader := csv.NewReader(e.config.Reader)

	for {
		lines, eofReached, err := getNextNLines(reader)
		if err != nil {
			wg.Wait()
			close(e.digestChannel)
			e.errorChannel <- err
			return
		}

		wg.Add(1)
		go e.digestForLines(lines, wg)

		if eofReached {
			break
		}
	}
	wg.Wait()
	close(e.digestChannel)
	e.errorChannel <- nil
}

func (e Engine) digestForLines(lines [][]string, wg *sync.WaitGroup) {
	output := make([]Digest, 0, len(lines))
	var createDigestFunc func(csv []string, pKey Positions, pRow Positions) Digest
	config := e.config

	if config.KeepSource {
		createDigestFunc = CreateDigestWithSource
	} else {
		createDigestFunc = CreateDigest
	}

	for _, line := range lines {
		output = append(output, createDigestFunc(line, config.Key, config.Value))
	}

	e.digestChannel <- output
	wg.Done()
}
