package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	t.Parallel()

	t.Run("EmptyEntries", func(t *testing.T) {
		randomRotation := &Random{}
		assert.NotNil(t, randomRotation)

		nextEntry := randomRotation.Next(nil)
		assert.Nil(t, nextEntry)
	})
}
