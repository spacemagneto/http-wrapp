package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightedRandom(t *testing.T) {
	t.Parallel()

	t.Run("EmptyEntries", func(t *testing.T) {
		weightedRandom := NewWeightedRandom()
		assert.NotNil(t, weightedRandom)

		nextEntry := weightedRandom.Select(nil)
		assert.Nil(t, nextEntry)
	})

	t.Run("ProxySelectionWithWeightPriorityCalculation", func(t *testing.T) {
		weightedRandom := NewWeightedRandom()
		assert.NotNil(t, weightedRandom)

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

	t.Run("SelectLastProxyWithMockRandFunc", func(t *testing.T) {
		weightedRandom := &WeightedSelector{
			randFloat64: func() float64 { return 1.1 },
		}

		assert.NotNil(t, weightedRandom)
		assert.NotNil(t, weightedRandom.randFloat64)

		entries := []*Entry{{proxy: &mockProxy{id: 1}}, {proxy: &mockProxy{id: 2}}}

		selected := weightedRandom.Select(entries)

		assert.NotNil(t, selected)
		assert.Equal(t, 2, selected.proxy.(*mockProxy).id, "Should return last entry when r exceeds sum")
	})
}
