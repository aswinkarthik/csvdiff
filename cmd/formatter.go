package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aswinkarthik93/csvdiff/pkg/digest"
)

// Formatter defines the interface through which differences
// can be formatted and displayed
type Formatter interface {
	Format(digest.Difference, io.Writer)
}

// RowMarkFormatter formats diff by marking each row as
// ADDED/MODIFIED. It mutates the row and adds as a new column.
type RowMarkFormatter struct{}

// Format prints the diff to os.Stdout
func (f *RowMarkFormatter) Format(diff digest.Difference, w io.Writer) {
	fmt.Fprintf(w, "Additions %d\n", len(diff.Additions))
	fmt.Fprintf(w, "Modifications %d\n", len(diff.Modifications))
	fmt.Fprintf(w, "Rows:\n")

	for _, added := range diff.Additions {
		fmt.Fprintf(w, "%s,%s\n", added, "ADDED")
	}

	for _, modified := range diff.Modifications {
		fmt.Fprintf(w, "%s,%s\n", modified, "MODIFIED")
	}
}

// JSONFormatter formats diff to as a JSON Object
type JSONFormatter struct{}

// Format prints the diff as a JSON
func (f *JSONFormatter) Format(diff digest.Difference, w io.Writer) {
	data, err := json.MarshalIndent(diff, "", "  ")

	if err != nil {
		panic(err)
	}

	w.Write(data)
}
