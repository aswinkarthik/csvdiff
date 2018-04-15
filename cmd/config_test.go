package cmd

import (
	"os"
	"testing"

	"github.com/aswinkarthik93/csv-digest/pkg/encoder"
	"github.com/stretchr/testify/assert"
)

func TestGetEncoder(t *testing.T) {
	config := Config{Encoder: "json"}
	assert.Equal(t, encoder.JsonEncoder{}, config.GetEncoder())

	config = Config{Encoder: "random"}
	assert.Equal(t, encoder.JsonEncoder{}, config.GetEncoder())
}

func TestGetKeyPositions(t *testing.T) {
	config := Config{KeyPositions: []int{0, 1}}
	assert.Equal(t, []int{0, 1}, config.GetKeyPositions())

	config = Config{KeyPositions: []int{}}
	assert.Equal(t, []int{0}, config.GetKeyPositions())
}

func TestReader(t *testing.T) {
	config := Config{Input: "STDIN"}
	assert.Equal(t, os.Stdin, config.GetReader())

	config = Config{Input: "-"}
	assert.Equal(t, os.Stdin, config.GetReader())
}

func TestWriter(t *testing.T) {
	config := Config{Input: "STDOUT"}
	assert.Equal(t, os.Stdout, config.GetWriter())

	config = Config{Input: "-"}
	assert.Equal(t, os.Stdout, config.GetWriter())
}
