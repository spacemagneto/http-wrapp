package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobin(t *testing.T) {
	t.Parallel()

	t.Run("EmptyEntries", func(t *testing.T) {
		roundRobin := &RoundRobin{}
		assert.NotNil(t, roundRobin)
		assert.Equal(t, uint64(0), roundRobin.counter.Load())

		nextProxy := roundRobin.Next(nil)
		assert.Nil(t, nextProxy)
		assert.Equal(t, uint64(0), roundRobin.counter.Load())
	})
}
