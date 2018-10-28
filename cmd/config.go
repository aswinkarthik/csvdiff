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
	PrimaryKeyPositions  []int
	ValueColumnPositions []int
	Format               string
}

// GetPrimaryKeys is to return the --primary-key flags as digest.Positions array.
func (c *Config) GetPrimaryKeys() digest.Positions {
	if len(c.PrimaryKeyPositions) > 0 {
		return c.PrimaryKeyPositions
	}
	return []int{0}
}

// GetValueColumns is to return the --value-columns flags as digest.Positions array.
func (c *Config) GetValueColumns() digest.Positions {
	if len(c.ValueColumnPositions) > 0 {
		return c.ValueColumnPositions
	}
	return []int{}
}

// Validate validates the config object
// and returns error if not valid.
func (c *Config) Validate() error {
	allFormats := []string{rowmark, jsonFormat}

	formatValid := false
	for _, format := range allFormats {
		if strings.ToLower(c.Format) == format {
			formatValid = true
		}
	}

	if !formatValid {
		return errors.New("Specified format is not valid")
	}

	return nil
}

const (
	rowmark    = "rowmark"
	jsonFormat = "json"
)

// Formatter instantiates a new formatted
// based on config.Format
func (c *Config) Formatter() Formatter {
	format := strings.ToLower(c.Format)
	if format == rowmark {
		return &RowMarkFormatter{}
	} else if format == jsonFormat {
		return &JSONFormatter{}
	}
	return &RowMarkFormatter{}
}
