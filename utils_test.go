package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopySet(t *testing.T) {
	cases := []struct {
		name string
		src  map[string]struct{}
		want map[string]struct{}
	}{
		{
			name: "NilSource",
			src:  nil,
			want: map[string]struct{}{},
		},
		{
			name: "EmptySource",
			src:  map[string]struct{}{},
			want: map[string]struct{}{},
		},
		{
			name: "SingleElement",
			src:  map[string]struct{}{"a": {}},
			want: map[string]struct{}{"a": {}},
		},
		{
			name: "MultipleElements",
			src:  map[string]struct{}{"a": {}, "b": {}, "c": {}},
			want: map[string]struct{}{"a": {}, "b": {}, "c": {}},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			expected := CopySet(tt.src)

			assert.Equal(t, tt.want, expected)
			assert.Equal(t, len(tt.want), len(expected))

			if tt.src != nil && len(tt.src) > 0 {
				for k := range tt.src {
					delete(tt.src, k)
					break
				}

				assert.NotEqual(t, len(tt.src), len(expected))
			}
		})
	}
}

func TestCopySetIntegerKeys(t *testing.T) {
	src := map[int]struct{}{1: {}, 2: {}, 100: {}}
	expected := CopySet(src)

	assert.Equal(t, src, expected)
	expected[999] = struct{}{}
	assert.NotContains(t, src, 999)
}
