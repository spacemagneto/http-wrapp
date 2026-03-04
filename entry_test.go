package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyEntry(t *testing.T) {
	t.Parallel()

	t.Run("InitialEntry", func(t *testing.T) {
		expectProxy := &mockProxy{id: 1}
		entry := newEntry(expectProxy)

		assert.NotNil(t, entry)
		assert.Equal(t, expectProxy, entry.Proxy())
		assert.NotNil(t, entry.Stats())
	})
}
