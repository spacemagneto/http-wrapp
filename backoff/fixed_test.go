package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFixedBackoffStrategy(t *testing.T) {
	t.Parallel()

	t.Run("InitStrategyWithoutDelay", func(t *testing.T) {
		fixed := NewFixed(0)
		assert.NotNil(t, fixed)
		assert.Equal(t, 1*time.Second, fixed.delay)
	})
}
