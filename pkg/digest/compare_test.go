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
100,col-1,col-2,col-3,hundred-value-modified
5,col-1,col-2,col-3,five-value-added
`

	baseConfig := &digest.Config{
		Reader: strings.NewReader(base),
		Key:    []int{0},
	}

	deltaConfig := &digest.Config{
		Reader: strings.NewReader(delta),
		Key:    []int{0},
	}

	expected := digest.Difference{
		Additions: []string{
			"4,col-1,col-2,col-3,four-value-added",
			"5,col-1,col-2,col-3,five-value-added",
		},
		Modifications: []string{
			"2,col-1,col-2,col-3,two-value-modified",
			"100,col-1,col-2,col-3,hundred-value-modified",
		},
	}

	actual := digest.Diff(baseConfig, deltaConfig)

	assert.ElementsMatch(t, expected.Modifications, actual.Modifications)
	assert.ElementsMatch(t, expected.Additions, actual.Additions)
}
