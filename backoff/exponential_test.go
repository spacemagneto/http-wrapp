package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialStrategy(t *testing.T) {
	t.Parallel()

	t.Run("InitStrategyWithDefaultStep", func(t *testing.T) {
		baseDelay := 1 * time.Second
		maxDelay := 10 * time.Second

		exponential := NewExponential(baseDelay, maxDelay)
		assert.NotNil(t, exponential)

		assert.Equal(t, baseDelay, exponential.baseDelay)
		assert.Equal(t, maxDelay, exponential.maxDelay)
		assert.Equal(t, float64(2), exponential.step)
	})
}
