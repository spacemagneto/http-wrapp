package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecorrelatedJitter(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitStrategy", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exponential := NewDecorrelatedJitter(baseDelay, maxDelay)
		assert.NotNil(t, exponential)

		assert.Equal(t, baseDelay, exponential.baseDelay)
		assert.Equal(t, maxDelay, exponential.maxDelay)
	})

	t.Run("InitStrategyWithDefaultValDelay", func(t *testing.T) {
		exponential := NewDecorrelatedJitter(0, 0)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultJitterBaseDelay, exponential.baseDelay)
		assert.Equal(t, DefaultJitterMaxDelay, exponential.maxDelay)
	})

	t.Run("InitStrategyWithInvalidArgs", func(t *testing.T) {
		exponential := NewDecorrelatedJitter(0, -10*time.Second)
		assert.NotNil(t, exponential)

		assert.Equal(t, DefaultJitterBaseDelay, exponential.baseDelay)
		assert.Equal(t, DefaultJitterMaxDelay, exponential.maxDelay)
	})

	t.Run("InitStrategyWhenMaxLessThanBase", func(t *testing.T) {
		baseDelay := 5 * time.Second
		maxDelay := 1 * time.Second
		jitter := NewDecorrelatedJitter(baseDelay, maxDelay)

		assert.Equal(t, baseDelay, jitter.maxDelay)
	})

	t.Run("SuccessGetNext", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		equalJitter := NewDecorrelatedJitter(baseDelay, maxDelay)
		delay := equalJitter.Next(1)
		assert.True(t, delay >= 0)
	})
}
