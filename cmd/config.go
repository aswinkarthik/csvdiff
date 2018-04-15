package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aswinkarthik93/csv-digest/pkg/encoder"
)

var encoders map[string]encoder.Encoder
var config Config

func init() {
	encoders = map[string]encoder.Encoder{"json": encoder.JsonEncoder{}}
	config = Config{}
}

type Config struct {
	KeyPositions  []int
	Encoder       string
	Base          string
	Delta         string
	Additions     string
	Modifications string
}

func (c Config) GetKeyPositions() []int {
	if len(c.KeyPositions) > 0 {
		return c.KeyPositions
	}
	return []int{0}
}

func (c Config) GetEncoder() encoder.Encoder {
	if val, ok := encoders[c.Encoder]; ok {
		return val
	} else {
		fmt.Println("Using JSON encoder")
		return encoders["json"]
	}
}

func (c Config) GetBaseReader() io.Reader {
	return getReader(c.Base)
}

func (c Config) GetDeltaReader() io.Reader {
	return getReader(c.Delta)
}

func (c Config) AdditionsWriter() io.WriteCloser {
	return getWriter(c.Additions)
}

func (c Config) ModificationsWriter() io.WriteCloser {
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

func GetEncoders() []string {
	result := make([]string, len(encoders))

	counter := 0
	for k := range encoders {
		result[counter] = k
		counter++
	}

	return result
}
