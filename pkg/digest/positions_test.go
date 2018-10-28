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
