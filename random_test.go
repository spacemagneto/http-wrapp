package client

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	t.Parallel()

	t.Run("EmptyEntries", func(t *testing.T) {
		randomRotation := &Random{}
		assert.NotNil(t, randomRotation)

		nextEntry := randomRotation.Select(nil)
		assert.Nil(t, nextEntry)
	})

	t.Run("SuccessSelectProxy", func(t *testing.T) {
		randomRotation := &Random{}
		assert.NotNil(t, randomRotation)

		randValue := rand.Int()
		expectProxy := &mockProxy{id: randValue}

		entries := []*Entry{newEntry(expectProxy)}
		nextProxy := randomRotation.Select(entries)
		assert.NotNil(t, nextProxy)
		assert.Equal(t, expectProxy, nextProxy.Proxy())
	})
}
