package digest_test

import (
	"strings"
	"testing"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/stretchr/testify/assert"
)

func TestPositionsMapValues(t *testing.T) {
	positions := digest.Positions([]int{0, 3})
	csv := []string{"zero", "one", "two", "three"}

	actual := positions.MapToValue(csv)
	expected := "zero,three"

	assert.Equal(t, expected, actual)
}

func TestPositionsMapValuesReturnsCompleteStringCsvIfEmpty(t *testing.T) {
	positions := digest.Positions([]int{})
	csv := []string{"zero", "one", "two", "three"}

	actual := positions.MapToValue(csv)
	expected := strings.Join(csv, digest.Separator)

	assert.Equal(t, expected, actual)
}

func TestPosition_Contains(t *testing.T) {
	positions := digest.Positions([]int{0, 3})

	assert.True(t, positions.Contains(3))
	assert.False(t, positions.Contains(4))
}

func TestPosition_Append(t *testing.T) {
	positions := digest.Positions([]int{0, 3})
	additionalPositions := digest.Positions([]int{4, 3})

	positions = positions.Append(additionalPositions)

	assert.ElementsMatch(t, []int{0, 3, 4}, []int(positions))
}
