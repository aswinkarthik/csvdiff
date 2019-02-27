package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
)

const (
	rowmark    = "rowmark"
	jsonFormat = "json"
	diffFormat = "diff"
)

// Formatter can print the differences to stdout
// and accompanying metadata to stderr
type Formatter struct {
	stdout io.Writer
	stderr io.Writer
	config Config
}

// NewFormatter can be used to create a new formatter
func NewFormatter(stdout, stderr io.Writer, config Config) *Formatter {
	return &Formatter{stdout: stdout, stderr: stderr, config: config}
}

// Format can be used to format the differences based on config
// to appropriate writers
func (f *Formatter) Format(diff digest.Differences) error {
	switch f.config.Format {
	case jsonFormat:
		return f.json(diff)
	case rowmark:
		return f.rowMark(diff)
	default:
		return fmt.Errorf("formatter not found")
	}
}

// JSONFormatter formats diff to as a JSON Object
// { "Additions": [...], "Modifications": [...] }
func (f *Formatter) json(diff digest.Differences) error {
	// jsonDifference is a struct to represent legacy JSON format
	type jsonDifference struct {
		Additions     []string
		Modifications []string
	}

	includes := config.GetIncludeColumnPositions()

	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, includes.MapToValue(addition))
	}

	modifications := make([]string, 0, len(diff.Modifications))
	for _, modification := range diff.Modifications {
		modifications = append(modifications, includes.MapToValue(modification.Current))
	}

	jsonDiff := jsonDifference{Additions: additions, Modifications: modifications}
	data, err := json.MarshalIndent(jsonDiff, "", "  ")

	if err != nil {
		return fmt.Errorf("error when serializing with JSON formatter: %v", err)
	}

	_, err = f.stdout.Write(data)

	if err != nil {
		return fmt.Errorf("error when writing to writer with JSON formatter: %v", err)
	}

	return nil
}

// RowMarkFormatter formats diff by marking each row as
// ADDED/MODIFIED. It mutates the row and adds as a new column.
func (f *Formatter) rowMark(diff digest.Differences) error {

	fmt.Fprintf(f.stderr, "Additions %d\n", len(diff.Additions))
	fmt.Fprintf(f.stderr, "Modifications %d\n", len(diff.Modifications))
	fmt.Fprintf(f.stderr, "Rows:\n")

	includes := config.GetIncludeColumnPositions()

	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, includes.MapToValue(addition))
	}

	modifications := make([]string, 0, len(diff.Modifications))
	for _, modification := range diff.Modifications {
		modifications = append(modifications, includes.MapToValue(modification.Current))
	}

	for _, added := range additions {
		_, err := fmt.Fprintf(f.stdout, "%s,%s\n", added, "ADDED")

		if err != nil {
			return fmt.Errorf("error when formatting additions with RowMark formatter: %v", err)
		}
	}

	for _, modified := range modifications {
		_, err := fmt.Fprintf(f.stdout, "%s,%s\n", modified, "MODIFIED")

		if err != nil {
			return fmt.Errorf("error when formatting modifications with RowMark formatter: %v", err)
		}

	}

	return nil
}
