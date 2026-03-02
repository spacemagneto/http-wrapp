package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedRandom(t *testing.T) {
	t.Parallel()

	t.Run("EmptyEntries", func(t *testing.T) {
		weighted := &WeightedRandom{}
		assert.NotNil(t, weighted)

		nextEntry := weighted.Next(nil)
		assert.Nil(t, nextEntry)
	})
}
