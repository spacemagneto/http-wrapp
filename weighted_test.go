package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedRandom(t *testing.T) {
	t.Parallel()

	weightedRandom := &WeightedRandom{}
	assert.NotNil(t, weightedRandom)

	t.Run("EmptyEntries", func(t *testing.T) {
		nextEntry := weightedRandom.Select(nil)
		assert.Nil(t, nextEntry)
	})

	t.Run("ProxySelectionWithWeightPriorityCalculation", func(t *testing.T) {
		first := &Entry{proxy: &mockProxy{id: 1}}
		first.stats.RecordSuccess()
		first.stats.RecordLatency(100)

		second := &Entry{proxy: &mockProxy{id: 2}}
		second.stats.RecordSuccess()
		second.stats.RecordLatency(500)

		entries := []*Entry{first, second}
		counts := make(map[int]int)
		iterations := 10000

		for i := 0; i < iterations; i++ {
			selected := weightedRandom.Select(entries)
			counts[selected.proxy.(*mockProxy).id]++
		}

		highestPriorityProxy := float64(counts[1]) / float64(iterations)
		otherPriorityProxy := float64(counts[2]) / float64(iterations)

		assert.InDelta(t, 0.83, highestPriorityProxy, 0.05)
		assert.InDelta(t, 0.17, otherPriorityProxy, 0.05)
	})
}
