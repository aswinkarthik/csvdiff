package cmd

import (
	"io"
	"log"
	"os"

	"github.com/aswinkarthik93/csvdiff/pkg/digest"
)

var config Config

func init() {
	config = Config{}
}

// Config is to store all command line Flags.
type Config struct {
	PrimaryKeyPositions  []int
	ValueColumnPositions []int
	Base                 string
	Delta                string
	Additions            string
	Modifications        string
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

// GetBaseReader returns an io.Reader for the base file.
func (c *Config) GetBaseReader() io.Reader {
	return getReader(c.Base)
}

// GetDeltaReader returns an io.Reader for the delta file.
func (c *Config) GetDeltaReader() io.Reader {
	return getReader(c.Delta)
}

// AdditionsWriter gives the output stream for the additions in delta csv.
func (c *Config) AdditionsWriter() io.WriteCloser {
	return getWriter(c.Additions)
}

// ModificationsWriter gives the output stream for the modifications in delta csv.
func (c *Config) ModificationsWriter() io.WriteCloser {
	return getWriter(c.Modifications)
}

func getReader(filename string) io.Reader {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	return file
}

func getWriter(outputStream string) io.WriteCloser {
	if outputStream != "STDOUT" {
		file, err := os.Create(outputStream)

		if err != nil {
			log.Fatal(err)
		}

		return file
	}
	return os.Stdout
}
