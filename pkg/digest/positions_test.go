package digest_test

import (
	"strings"
	"testing"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/stretchr/testify/assert"
)

const comma = ","

func TestPositions_MapValues(t *testing.T) {
	t.Run("should map positions to string", func(t *testing.T) {
		positions := digest.Positions([]int{0, 3})
		csv := []string{"zero", "one", "two", "three"}

		actual := positions.Join(csv, comma)
		expected := "zero,three"

		assert.Equal(t, expected, actual)
	})

	t.Run("should map all positions to string if positions is empty", func(t *testing.T) {
		positions := digest.Positions([]int{})
		csv := []string{"zero", "one", "two", "three"}

		actual := positions.Join(csv, comma)
		expected := strings.Join(csv, comma)

		assert.Equal(t, expected, actual)
	})

	t.Run("should not escape comma but retain new line if it is part of csv when mapping to values", func(t *testing.T) {
		positions := digest.Positions([]int{0, 3})
		csv := []string{"zero\n", "one", "two", "three,3"}

		actual := positions.Join(csv, comma)
		expected := "zero\n,three,3"

		assert.Equal(t, expected, actual)
	})
}

func TestPositions_String(t *testing.T) {
	t.Run("should map positions to string", func(t *testing.T) {
		positions := digest.Positions([]int{0, 3})
		csv := []string{"zero", "one", "two", "three"}

		actual := positions.String(csv, ',')
		expected := "zero,three"

		assert.Equal(t, expected, actual)
	})

	t.Run("should map positions to string using custom separator", func(t *testing.T) {
		positions := digest.Positions([]int{0, 3})
		csv := []string{"zero", "one", "two", "three"}

		actual := positions.String(csv, '|')
		expected := "zero|three"

		assert.Equal(t, expected, actual)
	})

	t.Run("should map all positions to string if positions is empty", func(t *testing.T) {
		positions := digest.Positions([]int{})
		csv := []string{"zero", "one", "two", "three"}

		actual := positions.String(csv, ',')
		expected := strings.Join(csv, comma)

		assert.Equal(t, expected, actual)
	})

	t.Run("should escape comma or new line if it is part of csv when mapping to values", func(t *testing.T) {
		positions := digest.Positions([]int{0, 3})
		csv := []string{"zero\n", "one", "two", "three,3"}

		actual := positions.String(csv, ',')
		expected := "\"zero\n\",\"three,3\""

		assert.Equal(t, expected, actual)
	})
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
