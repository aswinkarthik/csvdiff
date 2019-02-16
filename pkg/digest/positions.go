package digest

import (
	"strings"
)

// Positions represents positions of columns in a CSV array.
type Positions []int

// MapToValue plucks the values from CSV from
// their respective positions and concatenates
// them using Separator as a string.
func (p Positions) MapToValue(csv []string) string {
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
