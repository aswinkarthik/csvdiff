package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"strings"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
)

// Config is to store all command line Flags.
type Config struct {
	PrimaryKeyPositions    []int
	ValueColumnPositions   []int
	IncludeColumnPositions []int
	Format                 string
	BaseFilename           string
	DeltaFilename          string
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
func (c Config) GetIncludeColumnPositions() digest.Positions {
	if len(c.IncludeColumnPositions) > 0 {
		return c.IncludeColumnPositions
	}
	return c.GetValueColumns()
}

// Validate validates the config object
// and returns error if not valid.
func (c *Config) Validate(fs afero.Fs) error {
	{
		// format validation

		formatFound := false
		for _, format := range allFormats {
			if strings.ToLower(c.Format) == format {
				formatFound = true
			}
		}
		if !formatFound {
			return fmt.Errorf("specified format is not valid")
		}
	}

	{
		// base-file validation

		if exists, err := afero.Exists(fs, c.BaseFilename); err != nil {
			return fmt.Errorf("error reading base-file %s: %v", c.BaseFilename, err)
		} else if !exists {
			return fmt.Errorf("base-file %s does not exits", c.BaseFilename)
		}

		if isDir, err := afero.IsDir(fs, c.BaseFilename); err != nil {
			return fmt.Errorf("error reading base-file %s: %v", c.BaseFilename, err)
		} else if isDir {
			return fmt.Errorf("base-file %s should be a file", c.BaseFilename)
		}
	}

	{
		// delta file validation

		if exists, err := afero.Exists(fs, c.DeltaFilename); err != nil {
			return fmt.Errorf("error reading delta-file %s: %v", c.DeltaFilename, err)
		} else if !exists {
			return fmt.Errorf("delta-file %s does not exits", c.DeltaFilename)
		}

		if isDir, err := afero.IsDir(fs, c.DeltaFilename); err != nil {
			return fmt.Errorf("error reading delta-file %s: %v", c.DeltaFilename, err)
		} else if isDir {
			return fmt.Errorf("delta-file %s should be a file", c.DeltaFilename)
		}
	}

	return nil
}
