package digest_test

import (
	"strings"
	"testing"

	"github.com/aswinkarthik93/csvdiff/pkg/digest"
	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	baseDigest := map[uint64]uint64{
		10000106069522789940: 11608188164212916000,
		10000305084889337335: 11796412213504516000,
		10024909476616779194: 14500526491611670000,
		1004896778135186857:  15778011848259830000,
	}

	newDigest := map[uint64]uint64{
		10000106069522789940: 11608188164212916000,
		10000305084889337335: 11796412213504516001,
		10049141081086325814: 12259600610026582000,
	}

	additions, modifications := digest.Compare(baseDigest, newDigest)

	expectedAdditions := []uint64{10049141081086325814}
	expectedModifications := []uint64{10000305084889337335}

	assert.Equal(t, expectedAdditions, additions)
	assert.Equal(t, expectedModifications, modifications)
}

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
