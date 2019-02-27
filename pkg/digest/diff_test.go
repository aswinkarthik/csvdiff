package digest_test

import (
	"strings"
	"testing"

	"github.com/aswinkarthik/csvdiff/pkg/digest"
	"github.com/stretchr/testify/assert"
)

func TestDiff(t *testing.T) {
	base := `1,col-1,col-2,col-3,one-value
2,col-1,col-2,col-3,two-value
3,col-1,col-2,col-3,three-value
100,col-1,col-2,col-3,hundred-value
`

	delta := `1,col-1,col-2,col-3,one-value
2,col-1,col-2,col-3,two-value-modified
4,col-1,col-2,col-3,four-value-added
100,col-1-modified,col-2,col-3,hundred-value-modified
5,col-1,col-2,col-3,five-value-added
`

	t.Run("default config", func(t *testing.T) {
		baseConfig := &digest.Config{
			Reader:     strings.NewReader(base),
			Key:        []int{0},
			KeepSource: true,
		}

		deltaConfig := &digest.Config{
			Reader: strings.NewReader(delta),
			Key:    []int{0},
		}

		expected := digest.Differences{
			Additions: []digest.Addition{
				strings.Split("4,col-1,col-2,col-3,four-value-added", ","),
				strings.Split("5,col-1,col-2,col-3,five-value-added", ","),
			},
			Modifications: []digest.Modification{
				{
					Current:  strings.Split("2,col-1,col-2,col-3,two-value-modified", ","),
					Original: strings.Split("2,col-1,col-2,col-3,two-value", ","),
				},
				{
					Current:  strings.Split("100,col-1-modified,col-2,col-3,hundred-value-modified", ","),
					Original: strings.Split("100,col-1,col-2,col-3,hundred-value", ","),
				},
			},
		}

		actual, err := digest.Diff(*baseConfig, *deltaConfig)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("selective values columns without keeping source", func(t *testing.T) {
		baseConfig := &digest.Config{
			Reader:     strings.NewReader(base),
			Key:        []int{0},
			Value:      []int{1},
			KeepSource: false,
		}

		deltaConfig := &digest.Config{
			Reader: strings.NewReader(delta),
			Key:    []int{0},
			Value:  []int{1},
		}

		expected := digest.Differences{
			Additions: []digest.Addition{
				strings.Split("4,col-1,col-2,col-3,four-value-added", ","),
				strings.Split("5,col-1,col-2,col-3,five-value-added", ","),
			},
			Modifications: []digest.Modification{
				{
					Current: strings.Split("100,col-1-modified,col-2,col-3,hundred-value-modified", ","),
				},
			},
		}

		actual, err := digest.Diff(*baseConfig, *deltaConfig)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
