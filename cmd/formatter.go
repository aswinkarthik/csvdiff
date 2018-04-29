package cmd

import (
	"fmt"

	"github.com/aswinkarthik93/csvdiff/pkg/digest"
)

// Formatter defines the interface through which differences
// can be formatted and displayed
type Formatter interface {
	Format(digest.Difference)
}

// StdoutFormatter formats diff to STDOUT
type StdoutFormatter struct{}

// Format prints the diff to os.Stdout
func (f *StdoutFormatter) Format(diff digest.Difference) {
	fmt.Printf("Additions %d\n", len(diff.Additions))
	fmt.Printf("Modifications %d\n", len(diff.Modifications))
	fmt.Println("Rows:")

	for _, added := range diff.Additions {
		fmt.Printf("%s,%s\n", added, "ADDED")
	}

	for _, modified := range diff.Modifications {
		fmt.Printf("%s,%s\n", modified, "MODIFIED")
	}
}
