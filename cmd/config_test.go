package cmd

import (
	"testing"

	"github.com/aswinkarthik93/csvdiff/pkg/digest"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryKeyPositions(t *testing.T) {
	config := Config{PrimaryKeyPositions: []int{0, 1}}
	assert.Equal(t, digest.Positions([]int{0, 1}), config.GetPrimaryKeys())

	config = Config{PrimaryKeyPositions: []int{}}
	assert.Equal(t, digest.Positions([]int{0}), config.GetPrimaryKeys())

	config = Config{}
	assert.Equal(t, digest.Positions([]int{0}), config.GetPrimaryKeys())
}

func TestValueColumnPositions(t *testing.T) {
	config := Config{ValueColumnPositions: []int{0, 1}}
	assert.Equal(t, digest.Positions([]int{0, 1}), config.GetValueColumns())

	config = Config{ValueColumnPositions: []int{}}
	assert.Equal(t, digest.Positions([]int{}), config.GetValueColumns())

	config = Config{}
	assert.Equal(t, digest.Positions([]int{}), config.GetValueColumns())
}
