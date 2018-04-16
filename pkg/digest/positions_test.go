package digest

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPositionsMapValues(t *testing.T) {
	positions := Positions([]int{0, 3})
	csv := []string{"zero", "one", "two", "three"}

	actual := positions.MapToValue(csv)
	expected := "zero,three"

	assert.Equal(t, expected, actual)
}

func TestPositionsMapValuesReturnsCompleteStringCsvIfEmpty(t *testing.T) {
	positions := Positions([]int{})
	csv := []string{"zero", "one", "two", "three"}

	actual := positions.MapToValue(csv)
	expected := strings.Join(csv, Separator)

	assert.Equal(t, expected, actual)
}

func TestPositionsLength(t *testing.T) {
	positions := Positions([]int{0, 3})

	assert.Equal(t, 2, positions.Length())
}

func TestPositionsItems(t *testing.T) {
	items := []int{0, 3}
	positions := Positions(items)

	assert.Equal(t, items, positions.Items())
}
