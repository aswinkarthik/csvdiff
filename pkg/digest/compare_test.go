package digest

import (
	"testing"

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

	additions, modifications := Compare(baseDigest, newDigest)

	expectedAdditions := []uint64{10049141081086325814}
	expectedModifications := []uint64{10000305084889337335}

	assert.Equal(t, expectedAdditions, additions)
	assert.Equal(t, expectedModifications, modifications)
}
