package digest

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNextNLines(t *testing.T) {
	t.Run("should get given number of lines from csv", func(t *testing.T) {
		var csvBuilder strings.Builder
		const totalLines = 1000
		for i := 0; i < totalLines; i++ {
			csvBuilder.WriteString(fmt.Sprintf("%d,random-col-1,random-col-2\n", i))
		}

		csvFile := csv.NewReader(strings.NewReader(csvBuilder.String()))

		lines, eofReached, err := getNextNLines(csvFile)

		assert.Len(t, lines, bufferSize)
		assert.False(t, eofReached)
		assert.NoError(t, err)

		for i := 0; i < bufferSize; i++ {
			expected := []string{strconv.Itoa(i), "random-col-1", "random-col-2"}
			assert.Equal(t, expected, lines[i])
		}

		lines, eofReached, err = getNextNLines(csvFile)

		assert.Len(t, lines, totalLines-bufferSize)
		assert.True(t, eofReached)
		assert.NoError(t, err)

		for i := 0; i < totalLines-bufferSize; i++ {
			expected := []string{strconv.Itoa(i + bufferSize), "random-col-1", "random-col-2"}
			assert.Equal(t, expected, lines[i])
		}
	})

	t.Run("should throw error if not a valid csv", func(t *testing.T) {
		sampleInvalidCSV := `1,2,3
4,5,6
random-stuff
7,8,9`
		csvFile := csv.NewReader(strings.NewReader(sampleInvalidCSV))

		_, _, err := getNextNLines(csvFile)

		assert.Error(t, err)
	})
}
