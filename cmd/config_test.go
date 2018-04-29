package cmd_test

import (
	"testing"

	"github.com/aswinkarthik93/csvdiff/cmd"
	"github.com/aswinkarthik93/csvdiff/pkg/digest"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryKeyPositions(t *testing.T) {
	config := cmd.Config{PrimaryKeyPositions: []int{0, 1}}
	assert.Equal(t, digest.Positions([]int{0, 1}), config.GetPrimaryKeys())

	config = cmd.Config{PrimaryKeyPositions: []int{}}
	assert.Equal(t, digest.Positions([]int{0}), config.GetPrimaryKeys())

	config = cmd.Config{}
	assert.Equal(t, digest.Positions([]int{0}), config.GetPrimaryKeys())
}

func TestValueColumnPositions(t *testing.T) {
	config := cmd.Config{ValueColumnPositions: []int{0, 1}}
	assert.Equal(t, digest.Positions([]int{0, 1}), config.GetValueColumns())

	config = cmd.Config{ValueColumnPositions: []int{}}
	assert.Equal(t, digest.Positions([]int{}), config.GetValueColumns())

	config = cmd.Config{}
	assert.Equal(t, digest.Positions([]int{}), config.GetValueColumns())
}

func TestConfigValidate(t *testing.T) {
	var config *cmd.Config

	config = &cmd.Config{}
	assert.Error(t, config.Validate())

	config = &cmd.Config{Format: "rowmark"}
	assert.NoError(t, config.Validate())

	config = &cmd.Config{Format: "rowMARK"}
	assert.NoError(t, config.Validate())

	config = &cmd.Config{Format: "json"}
	assert.NoError(t, config.Validate())
}

func TestDefaultConfigFormatter(t *testing.T) {
	config := &cmd.Config{}

	formatter, ok := config.Formatter().(*cmd.RowMarkFormatter)

	assert.True(t, ok)
	assert.NotNil(t, formatter)
}

func TestConfigFormatter(t *testing.T) {
	var config *cmd.Config
	var formatter cmd.Formatter
	var ok bool

	config = &cmd.Config{Format: "rowmark"}
	formatter, ok = config.Formatter().(*cmd.RowMarkFormatter)
	assert.True(t, ok)
	assert.NotNil(t, formatter)

	config = &cmd.Config{Format: "json"}
	formatter, ok = config.Formatter().(*cmd.JSONFormatter)
	assert.True(t, ok)
	assert.NotNil(t, formatter)
}
