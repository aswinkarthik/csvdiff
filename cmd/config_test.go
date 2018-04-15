package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyPositions(t *testing.T) {
	config := Config{KeyPositions: []int{0, 1}}
	assert.Equal(t, []int{0, 1}, config.GetKeyPositions())

	config = Config{KeyPositions: []int{}}
	assert.Equal(t, []int{0}, config.GetKeyPositions())
}
