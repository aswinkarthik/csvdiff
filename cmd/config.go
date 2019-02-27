package cmd

import (
	"errors"
	"strings"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
)

var config Config

func init() {
	config = Config{}
}

// Config is to store all command line Flags.
type Config struct {
	PrimaryKeyPositions    []int
	ValueColumnPositions   []int
	IncludeColumnPositions []int
	Format                 string
}

// GetPrimaryKeys is to return the --primary-key flags as digest.Positions array.
func (c *Config) GetPrimaryKeys() digest.Positions {
	if len(c.PrimaryKeyPositions) > 0 {
		return c.PrimaryKeyPositions
	}
	return []int{0}
}

// GetValueColumns is to return the --columns flags as digest.Positions array.
func (c *Config) GetValueColumns() digest.Positions {
	if len(c.ValueColumnPositions) > 0 {
		return c.ValueColumnPositions
	}
	return []int{}
}

// GetIncludeColumnPositions is to return the --include flags as digest.Positions array.
// If empty, it is value columns
func (c *Config) GetIncludeColumnPositions() digest.Positions {
	if len(c.IncludeColumnPositions) > 0 {
		return c.IncludeColumnPositions
	}
	return c.GetValueColumns()
}

// Validate validates the config object
// and returns error if not valid.
func (c *Config) Validate() error {

	for _, format := range allFormats {
		if strings.ToLower(c.Format) == format {
			return nil
		}
	}

	return errors.New("Specified format is not valid")
}
