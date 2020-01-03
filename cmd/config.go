package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/afero"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
)

// Context is to store all command line Flags.
type Context struct {
	fs                     afero.Fs
	primaryKeyPositions    []int
	valueColumnPositions   []int
	includeColumnPositions []int
	format                 string
	baseFilename           string
	deltaFilename          string
	baseFile               afero.File
	deltaFile              afero.File
	recordCount            int
	separator              rune
	lazyQuotes             bool
}

// NewContext can take all CLI flags and create a cmd.Context
// Validations are done as part of this.
// File pointers are created too.
func NewContext(
	fs afero.Fs,
	primaryKeyPositions []int,
	valueColumnPositions []int,
	ignoreValueColumnPositions []int,
	includeColumnPositions []int,
	format string,
	baseFilename string,
	deltaFilename string,
	separator rune,
	lazyQuotes bool,
) (*Context, error) {
	baseRecordCount, err := getColumnsCount(fs, baseFilename, separator, lazyQuotes)
	if err != nil {
		return nil, fmt.Errorf("error in base-file: %v", err)
	}

	deltaRecordCount, err := getColumnsCount(fs, deltaFilename, separator, lazyQuotes)
	if err != nil {
		return nil, fmt.Errorf("error in delta-file: %v", err)
	}

	if baseRecordCount != deltaRecordCount {
		return nil, fmt.Errorf("base-file and delta-file columns count do not match")
	}

	if len(ignoreValueColumnPositions) > 0 && len(valueColumnPositions) > 0 {
		return nil, fmt.Errorf("only one of --columns or --ignore-columns")
	}
	if len(ignoreValueColumnPositions) > 0 {
		valueColumnPositions = inferValueColumns(baseRecordCount, ignoreValueColumnPositions)
	}

	baseFile, err := fs.Open(baseFilename)
	if err != nil {
		return nil, err
	}
	deltaFile, err := fs.Open(deltaFilename)
	if err != nil {
		return nil, err
	}
	ctx := &Context{
		fs:                     fs,
		primaryKeyPositions:    primaryKeyPositions,
		valueColumnPositions:   valueColumnPositions,
		includeColumnPositions: includeColumnPositions,
		format:                 format,
		baseFilename:           baseFilename,
		deltaFilename:          deltaFilename,
		baseFile:               baseFile,
		deltaFile:              deltaFile,
		recordCount:            baseRecordCount,
		separator:              separator,
		lazyQuotes:             lazyQuotes,
	}

	if err := ctx.validate(); err != nil {
		return nil, fmt.Errorf("validation failed: %v", err)
	}

	return ctx, nil
}

// GetPrimaryKeys is to return the --primary-key flags as digest.Positions array.
func (c *Context) GetPrimaryKeys() digest.Positions {
	if len(c.primaryKeyPositions) > 0 {
		return c.primaryKeyPositions
	}
	return []int{0}
}

// GetValueColumns is to return the --columns flags as digest.Positions array.
func (c *Context) GetValueColumns() digest.Positions {
	if len(c.valueColumnPositions) > 0 {
		return c.valueColumnPositions
	}
	return []int{}
}

// GetIncludeColumnPositions is to return the --include flags as digest.Positions array.
// If empty, it is value columns
func (c Context) GetIncludeColumnPositions() digest.Positions {
	if len(c.includeColumnPositions) > 0 {
		return c.includeColumnPositions
	}
	return c.GetValueColumns()
}

// validate validates the context object
// and returns error if not valid.
func (c *Context) validate() error {
	{
		// format validation

		formatFound := false
		for _, format := range allFormats {
			if strings.ToLower(c.format) == format {
				formatFound = true
			}
		}
		if !formatFound {
			return fmt.Errorf("specified format is not valid")
		}
	}

	{
		comparator := func(element int) bool {
			return element < c.recordCount
		}

		if !assertAll(c.primaryKeyPositions, comparator) {
			return fmt.Errorf("--primary-key positions are out of bounds")
		}
		if !assertAll(c.includeColumnPositions, comparator) {
			return fmt.Errorf("--include positions are out of bounds")
		}
		if !assertAll(c.valueColumnPositions, comparator) {
			return fmt.Errorf("--columns positions are out of bounds")
		}
	}

	return nil
}

func inferValueColumns(recordCount int, ignoreValueColumns []int) digest.Positions {
	lookupMap := make(map[int]struct{})
	for _, pos := range ignoreValueColumns {
		lookupMap[pos] = struct{}{}
	}

	valueColumns := make(digest.Positions, 0)
	if len(ignoreValueColumns) > 0 {
		for i := 0; i < recordCount; i++ {
			if _, exists := lookupMap[i]; !exists {
				valueColumns = append(valueColumns, i)
			}
		}
	}

	return valueColumns
}

func assertAll(elements []int, assertFn func(element int) bool) bool {
	for _, el := range elements {
		if !assertFn(el) {
			return false
		}
	}
	return true
}

func getColumnsCount(fs afero.Fs, filename string, separator rune, lazyQuotes bool) (int, error) {
	base, err := fs.Open(filename)
	if err != nil {
		return 0, err
	}
	defer base.Close()
	csvReader := csv.NewReader(base)
	csvReader.Comma = separator
	csvReader.LazyQuotes = lazyQuotes
	record, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return 0, fmt.Errorf("unable to process headers from csv file. EOF reached. invalid CSV file")
		}
		return 0, err
	}

	return len(record), nil
}

// BaseDigestConfig creates a digest.Context from cmd.Context
// that is needed to start the diff process
func (c *Context) BaseDigestConfig() (digest.Config, error) {
	return digest.Config{
		Reader:     c.baseFile,
		Value:      c.valueColumnPositions,
		Key:        c.primaryKeyPositions,
		Include:    c.includeColumnPositions,
		Separator:  c.separator,
		LazyQuotes: c.lazyQuotes,
	}, nil
}

// DeltaDigestConfig creates a digest.Context from cmd.Context
// that is needed to start the diff process
func (c *Context) DeltaDigestConfig() (digest.Config, error) {
	return digest.Config{
		Reader:     c.deltaFile,
		Value:      c.valueColumnPositions,
		Key:        c.primaryKeyPositions,
		Include:    c.includeColumnPositions,
		Separator:  c.separator,
		LazyQuotes: c.lazyQuotes,
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
