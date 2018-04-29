package digest

import "strings"

// Positions represents positions of columns in a CSV array.
type Positions []int

// MapToValue plucks the values from CSV from
// their respective positions and concatenates
// them using Separator as a string.
func (p Positions) MapToValue(csv []string) string {
	if len(p) == 0 {
		return strings.Join(csv, Separator)
	}
	output := make([]string, len(p))
	for i, pos := range p {
		output[i] = csv[pos]
	}
	return strings.Join(output, Separator)
}
