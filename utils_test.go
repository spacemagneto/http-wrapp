package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopySetIntegerKeys(t *testing.T) {
	src := map[int]struct{}{1: {}, 2: {}, 100: {}}
	expected := CopySet(src)

	assert.Equal(t, src, expected)
	expected[999] = struct{}{}
	assert.NotContains(t, src, 999)
}
