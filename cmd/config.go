package cmd

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"strings"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
)

// Context is to store all command line Flags.
type Context struct {
	PrimaryKeyPositions    []int
	ValueColumnPositions   []int
	IncludeColumnPositions []int
	Format                 string
	BaseFilename           string
	DeltaFilename          string
	baseFile               afero.File
	deltaFile              afero.File
}

// GetPrimaryKeys is to return the --primary-key flags as digest.Positions array.
func (c *Context) GetPrimaryKeys() digest.Positions {
	if len(c.PrimaryKeyPositions) > 0 {
		return c.PrimaryKeyPositions
	}
	return []int{0}
}

// GetValueColumns is to return the --columns flags as digest.Positions array.
func (c *Context) GetValueColumns() digest.Positions {
	if len(c.ValueColumnPositions) > 0 {
		return c.ValueColumnPositions
	}
	return []int{}
}

// GetIncludeColumnPositions is to return the --include flags as digest.Positions array.
// If empty, it is value columns
func (c Context) GetIncludeColumnPositions() digest.Positions {
	if len(c.IncludeColumnPositions) > 0 {
		return c.IncludeColumnPositions
	}
	return c.GetValueColumns()
}

// Validate validates the context object
// and returns error if not valid.
func (c *Context) Validate(fs afero.Fs) error {
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

	{
		baseRecordCount, err := getColumnsCount(fs, c.BaseFilename)
		if err != nil {
			return err
		}
		deltaRecordCount, err := getColumnsCount(fs, c.DeltaFilename)
		if err != nil {
			return err
		}

		if baseRecordCount != deltaRecordCount {
			return fmt.Errorf("base-file and delta-file columns count do not match")
		}

		comparator := func(element int) bool {
			return element < baseRecordCount
		}

		if !assertAll(c.PrimaryKeyPositions, comparator) {
			return fmt.Errorf("--primary-key positions are out of bounds")
		}
		if !assertAll(c.IncludeColumnPositions, comparator) {
			return fmt.Errorf("--include positions are out of bounds")
		}
		if !assertAll(c.ValueColumnPositions, comparator) {
			return fmt.Errorf("--columns positions are out of bounds")
		}
	}

	return nil
}

func assertAll(elements []int, assertFn func(element int) bool) bool {
	for _, el := range elements {
		if !assertFn(el) {
			return false
		}
	}
	return true
}

func getColumnsCount(fs afero.Fs, filename string) (int, error) {
	base, err := fs.Open(filename)
	if err != nil {
		return 0, err
	}
	defer base.Close()
	csvReader := csv.NewReader(base)

	record, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("unable to process headers from csv file. EOF reached")
		}
		return 0, err
	}

	return len(record), nil
}

// BaseDigestConfig creates a digest.Context from cmd.Context
// that is needed to start the diff process
func (c *Context) BaseDigestConfig(fs afero.Fs) (digest.Config, error) {
	baseFile, err := fs.Open(c.BaseFilename)
	if err != nil {
		return digest.Config{}, err
	}

	c.baseFile = baseFile

	return digest.Config{
		Reader:  baseFile,
		Value:   c.ValueColumnPositions,
		Key:     c.PrimaryKeyPositions,
		Include: c.IncludeColumnPositions,
	}, nil
}

// DeltaDigestConfig creates a digest.Context from cmd.Context
// that is needed to start the diff process
func (c *Context) DeltaDigestConfig(fs afero.Fs) (digest.Config, error) {
	deltaFile, err := fs.Open(c.DeltaFilename)
	if err != nil {
		return digest.Config{}, err
	}

	c.baseFile = deltaFile

	return digest.Config{
		Reader:  deltaFile,
		Value:   c.ValueColumnPositions,
		Key:     c.PrimaryKeyPositions,
		Include: c.IncludeColumnPositions,
	}, nil
}

// Close all file handles
func (c *Context) Close() {
	if c.baseFile != nil {
		_ = c.baseFile.Close()
	}
	if c.deltaFile != nil {
		_ = c.deltaFile.Close()
	}
}
