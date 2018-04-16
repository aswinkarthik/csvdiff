package digest

import "strings"

type Positions []int

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

func (p Positions) Length() int {
	return len([]int(p))
}

func (p Positions) Items() []int {
	return []int(p)
}
