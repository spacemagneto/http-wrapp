package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProxyEntry(t *testing.T) {
	t.Parallel()

	t.Run("InitialEntry", func(t *testing.T) {
		expectProxy := &mockProxy{id: 1}
		entry := newEntry(expectProxy)

		assert.NotNil(t, entry)
		assert.Equal(t, expectProxy, entry.Proxy())
		assert.NotNil(t, entry.Stats())
	})

	t.Run("HealthCheckHealthy", func(t *testing.T) {
		expectProxy := &mockProxy{id: 1}
		entry := newEntry(expectProxy)

		assert.NotNil(t, entry)
		assert.Equal(t, expectProxy, entry.Proxy())
		assert.NotNil(t, entry.Stats())

		isHealthy := entry.HealthCheck(3, time.Minute)
		assert.True(t, isHealthy)
	})

	t.Run("HealthCheckWithQuarantine", func(t *testing.T) {
		expectProxy := &mockProxy{id: 1}
		entry := newEntry(expectProxy)

		entry.stats.RecordFailed()
		entry.stats.RecordFailed()
		entry.stats.RecordFailed()

		isHealthy := entry.HealthCheck(3, time.Minute)
		assert.False(t, isHealthy)
	})
}
