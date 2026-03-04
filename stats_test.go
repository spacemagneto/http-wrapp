package client

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	t.Parallel()

	t.Run("InitialState", func(t *testing.T) {
		stats := &Stats{}
		assert.Equal(t, int64(0), stats.ConsecutiveFails())
		assert.Equal(t, int64(0), stats.SuccessCount())
		assert.Equal(t, int64(0), stats.Failures())
		assert.Equal(t, 0.0, stats.AvgLatencyMs())
		assert.Equal(t, baseWeight, stats.successRate())
		assert.Equal(t, baseWeight, stats.Weight())
		assert.True(t, stats.LastFailedTime().IsZero())
		assert.True(t, stats.LastUsedTime().IsZero())
	})

	t.Run("RecordSuccess", func(t *testing.T) {
		stats := &Stats{}

		stats.RecordFailed()
		stats.RecordFailed()
		assert.Equal(t, int64(2), stats.ConsecutiveFails())

		stats.RecordSuccess()
		assert.Equal(t, int64(0), stats.ConsecutiveFails())
		assert.Equal(t, int64(1), stats.SuccessCount())
		assert.Equal(t, int64(2), stats.Failures())
	})

	t.Run("RecordFailed", func(t *testing.T) {
		stats := &Stats{}

		beforeFailed := time.Now().UnixNano()
		stats.RecordFailed()
		afterFailed := time.Now().UnixNano()

		assert.Equal(t, int64(1), stats.ConsecutiveFails())
		assert.Equal(t, int64(1), stats.Failures())

		lastFailed := stats.LastFailedTime().UnixNano()
		assert.True(t, lastFailed >= beforeFailed && lastFailed <= afterFailed)
	})

	t.Run("RecordLatency", func(t *testing.T) {
		stats := &Stats{}
		stats.RecordLatency(50)
		stats.RecordLatency(50)
		stats.RecordLatency(300)
		stats.RecordLatency(200)

		assert.Equal(t, float64(150), stats.AvgLatencyMs())
	})

	t.Run("LastUsed", func(t *testing.T) {
		stats := &Stats{}

		beforeUsed := time.Now().UnixNano()
		stats.LastUsed()
		afterUsed := time.Now().UnixNano()

		lastUsed := stats.LastUsedTime().UnixNano()
		assert.True(t, lastUsed >= beforeUsed && lastUsed <= afterUsed)
	})

	t.Run("SuccessRate", func(t *testing.T) {
		stats := &Stats{}

		stats.RecordSuccess()
		stats.RecordSuccess()
		stats.RecordFailed()

		proxyRate := stats.successRate()
		assert.Equal(t, 0.6, math.Trunc(proxyRate*10)/10)
	})
}
