package cmd

import (
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
