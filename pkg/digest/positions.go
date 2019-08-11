package digest

import (
	csvlib "encoding/csv"
	"strings"
)

// Positions represents positions of columns in a CSV array.
type Positions []int

// Join plucks the values from CSV from
// their respective positions and concatenates
// them using Separator as a string.
func (p Positions) Join(csv []string) string {
	if len(p) == 0 {
		return strings.Join(csv, Separator)
	}

	csvStr := strings.Builder{}
	for _, pos := range p[:len(p)-1] {
		csvStr.WriteString(csv[pos])
		csvStr.WriteString(Separator)
	}
	csvStr.WriteString(csv[p[len(p)-1]])
	return csvStr.String()
}

// String method converts to csv mapping to positions
// escapes necessary characters
func (p Positions) String(csv []string) string {
	selectiveCsv := csv
	if len(p) != 0 {
		selectiveCsv = make([]string, 0, len(p))
		for _, pos := range p {
			selectiveCsv = append(selectiveCsv, csv[pos])
		}
	}

	csvStr := strings.Builder{}
	w := csvlib.NewWriter(&csvStr)
	_ = w.Write(selectiveCsv)
	w.Flush()
	csvWithNewLine := csvStr.String()
	return csvWithNewLine[:len(csvWithNewLine)-1]
}

// Append additional positions to existing positions.
// Imp: Removes Duplicate. Does not mutate the original array
func (p Positions) Append(additional Positions) Positions {
	for _, toBeAdded := range additional {
		if !p.Contains(toBeAdded) {
			p = append(p, toBeAdded)
		}
	}

	return p
}

// Contains returns true if position is already present in Positions
func (p Positions) Contains(position int) bool {
	for _, each := range p {
		if each == position {
			return true
		}
	}

	return false
}
