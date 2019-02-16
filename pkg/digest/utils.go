package digest

import (
	"encoding/csv"
	"io"
)

func getNextNLines(reader *csv.Reader) ([][]string, bool, error) {
	lines := make([][]string, bufferSize)

	lineCount := 0
	eofReached := false
	for ; lineCount < bufferSize; lineCount++ {
		line, err := reader.Read()
		lines[lineCount] = line
		if err != nil {
			if err == io.EOF {
				eofReached = true
				break
			}

			return nil, true, err
		}
	}

	return lines[:lineCount], eofReached, nil
}
