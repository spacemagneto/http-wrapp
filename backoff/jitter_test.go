package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJitterStrategy(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exponential := NewJitter(baseDelay, maxDelay)
		assert.NotNil(t, exponential)

		assert.Equal(t, baseDelay, exponential.baseDelay)
		assert.Equal(t, maxDelay, exponential.maxDelay)
	})

	t.Run("InitStrategyWithDefaultValDelay", func(t *testing.T) {
		exponential := NewJitter(0, 0)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultJitterBaseDelay, exponential.baseDelay)
		assert.Equal(t, DefaultJitterMaxDelay, exponential.maxDelay)
	})

	t.Run("InitStrategyWithInvalidDelay", func(t *testing.T) {
		exponential := NewJitter(0, -10*time.Second)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultJitterBaseDelay, exponential.baseDelay)
		assert.Equal(t, DefaultJitterMaxDelay, exponential.maxDelay)
	})
}
