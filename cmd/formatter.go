package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
)

// Formatter defines the interface through which differences
// can be formatted and displayed
type Formatter interface {
	Format(digest.Differences) error
}

// RowMarkFormatter formats diff by marking each row as
// ADDED/MODIFIED. It mutates the row and adds as a new column.
type RowMarkFormatter struct {
	Stdout io.Writer
	Stderr io.Writer
}

// Format prints the diff to os.Stdout
func (f *RowMarkFormatter) Format(diff digest.Differences) error {
	fmt.Fprintf(f.Stderr, "Additions %d\n", len(diff.Additions))
	fmt.Fprintf(f.Stderr, "Modifications %d\n", len(diff.Modifications))
	fmt.Fprintf(f.Stderr, "Rows:\n")
	
	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, strings.Join(addition, ","))
	}

	modifications := make([]string, 0, len(diff.Modifications))
	for _, modification := range diff.Modifications {
		modifications = append(modifications, strings.Join(modification.Current, ","))
	}

	for _, added := range additions {
		_, err := fmt.Fprintf(f.Stdout, "%s,%s\n", added, "ADDED")

		if err != nil {
			return fmt.Errorf("error when formatting additions with RowMark formatter: %v", err)
		}
	}

	for _, modified := range modifications {
		_, err := fmt.Fprintf(f.Stdout, "%s,%s\n", modified, "MODIFIED")

		if err != nil {
			return fmt.Errorf("error when formatting modifications with RowMark formatter: %v", err)
		}

	}

	return nil
}

// JSONFormatter formats diff to as a JSON Object
type JSONFormatter struct {
	Stdout io.Writer
}

// JSONDifference is a struct to represent legacy JSON format
type JSONDifference struct {
	Additions     []string
	Modifications []string
}

// Format prints the diff as a JSON
func (f *JSONFormatter) Format(diff digest.Differences) error {
	additions := make([]string, 0, len(diff.Additions))
	for _, addition := range diff.Additions {
		additions = append(additions, strings.Join(addition, ","))
	}

	modifications := make([]string, 0, len(diff.Modifications))
	for _, modification := range diff.Modifications {
		modifications = append(modifications, strings.Join(modification.Current, ","))
	}

	jsonDiff := JSONDifference{Additions: additions, Modifications: modifications}
	data, err := json.MarshalIndent(jsonDiff, "", "  ")

	if err != nil {
		return fmt.Errorf("error when serializing with JSON formatter: %v", err)
	}

	_, err = f.Stdout.Write(data)

	if err != nil {
		return fmt.Errorf("error when writing to writer with JSON formatter: %v", err)
	}

	return nil
}
