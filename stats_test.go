package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	t.Parallel()

	t.Run("InitialState", func(t *testing.T) {
		stats := &Stats{}
		assert.Equal(t, int32(0), stats.ConsecutiveFails())
		assert.Equal(t, int64(0), stats.SuccessCount())
		assert.Equal(t, int64(0), stats.TotalFails())
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
		assert.Equal(t, int32(2), stats.ConsecutiveFails())

		stats.RecordSuccess()
		assert.Equal(t, int32(0), stats.ConsecutiveFails())
		assert.Equal(t, int64(1), stats.SuccessCount())
		assert.Equal(t, int64(2), stats.TotalFails())
	})
}
