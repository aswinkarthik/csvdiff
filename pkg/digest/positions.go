package digest

import "strings"

// Positions represents positions of columns in a CSV array.
type Positions []int

// MapToValue plucks the values from CSV from
// their respective positions and concatenates
// them using Separator as a string.
func (p Positions) MapToValue(csv []string) string {
	if p.Length() == 0 {
		return strings.Join(csv, Separator)
	}
	output := make([]string, p.Length())
	for i, pos := range p.Items() {
		output[i] = csv[pos]
	}
	return strings.Join(output, Separator)
}

// Length returns the size of the Positions array.
func (p Positions) Length() int {
	return len([]int(p))
}

// Items returns the elements of the Positions array
// as an array of int
func (p Positions) Items() []int {
	return []int(p)
}
