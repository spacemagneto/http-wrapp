package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProxyPool(t *testing.T) {
	t.Parallel()

	t.Run("DefaultPoolConfigValues", func(t *testing.T) {
		cfg := defaultPoolConfig()

		assert.Equal(t, int64(3), cfg.MaxFails)
		assert.Equal(t, 30*time.Second, cfg.CooldownWindow)
		assert.NotNil(t, cfg.Selector)

		_, ok := cfg.Selector.(*RoundRobinSelector)
		assert.True(t, ok)
	})
}
