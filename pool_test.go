package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProxyPool(t *testing.T) {
	t.Parallel()

	t.Run("PoolConfigCustomValues", func(t *testing.T) {
		selector := &RoundRobinSelector{}
		cfg := PoolConfig{Selector: selector}

		pool := NewPool([]Proxy{&mockProxy{}}, cfg)

		assert.Equal(t, int64(3), pool.cfg.MaxFails)
		assert.Equal(t, 30*time.Second, pool.cfg.CooldownWindow)
		assert.Equal(t, selector, pool.cfg.Selector)
	})

	t.Run("DefaultPoolConfigValues", func(t *testing.T) {
		cfg := defaultPoolConfig()

		assert.Equal(t, int64(3), cfg.MaxFails)
		assert.Equal(t, 30*time.Second, cfg.CooldownWindow)
		assert.NotNil(t, cfg.Selector)

		_, ok := cfg.Selector.(*RoundRobinSelector)
		assert.True(t, ok)
	})

	t.Run("HealthyEntriesFilteringLogic", func(t *testing.T) {
		proxies := []Proxy{&mockProxy{id: 1}, &mockProxy{id: 2}, &mockProxy{id: 3}}

		cfg := PoolConfig{MaxFails: 2, CooldownWindow: time.Minute}

		pool := NewPool(proxies, cfg)
		assert.NotNil(t, pool)

		pool.entries[1].Stats().RecordFailed()
		pool.entries[2].Stats().RecordFailed()
		pool.entries[2].Stats().RecordFailed()
		pool.entries[2].Stats().RecordFailed()

		healthy := pool.healthyEntries()

		assert.Len(t, healthy, 2)
		assert.Equal(t, pool.entries[0], healthy[0])
		assert.Equal(t, pool.entries[1], healthy[1])
	})

	t.Run("EmptyHealthyEntriesWithAllProxyInQuarantine", func(t *testing.T) {
		proxies := []Proxy{&mockProxy{id: 1}, &mockProxy{id: 2}}

		cfg := PoolConfig{MaxFails: 2, CooldownWindow: time.Minute}

		pool := NewPool(proxies, cfg)
		assert.NotNil(t, pool)

		pool.entries[0].Stats().RecordFailed()
		pool.entries[0].Stats().RecordFailed()

		pool.entries[1].Stats().RecordFailed()
		pool.entries[1].Stats().RecordFailed()

		healthy := pool.healthyEntries()

		assert.Empty(t, healthy)
	})

	t.Run("SuccessPickEntry", func(t *testing.T) {
		proxies := []Proxy{&mockProxy{id: 1}, &mockProxy{id: 2}}

		selector := &RoundRobinSelector{}
		cfg := PoolConfig{MaxFails: 2, CooldownWindow: time.Minute, Selector: selector}

		pool := NewPool(proxies, cfg)
		assert.NotNil(t, pool)

		entry, err := pool.Pick()
		assert.NoError(t, err)
		assert.NotNil(t, entry)

		assert.Equal(t, uint64(1), selector.counter.Load())
	})

	t.Run("PoolWithEmptyProxies", func(t *testing.T) {
		var proxies []Proxy

		selector := &RoundRobinSelector{}
		cfg := PoolConfig{MaxFails: 2, CooldownWindow: time.Minute, Selector: selector}

		pool := NewPool(proxies, cfg)
		assert.NotNil(t, pool)

		_, err := pool.Pick()
		assert.Error(t, err)
		assert.ErrorIs(t, ErrProxyPoolEmpty, err)
	})

	t.Run("PickStaleProxyBecauseNoNaveHealthyEntries", func(t *testing.T) {
		proxies := []Proxy{&mockProxy{id: 1}, &mockProxy{id: 2}}

		selector := &RoundRobinSelector{}
		cfg := PoolConfig{MaxFails: 2, CooldownWindow: time.Minute, Selector: selector}

		pool := NewPool(proxies, cfg)
		assert.NotNil(t, pool)

		pool.entries[0].Stats().RecordFailed()
		pool.entries[0].Stats().RecordFailed()

		pool.entries[1].Stats().RecordFailed()
		pool.entries[1].Stats().RecordFailed()

		entry, err := pool.Pick()
		assert.NoError(t, err)
		assert.NotNil(t, entry)

		assert.Equal(t, uint64(1), selector.counter.Load())
	})
}
