package client

import "math/rand/v2"

// WeightedSelector picks an entry at random, with probability proportional
// to each entry's Weight (derived from success rate and average latency).
//
// Unlike a deterministic "pick the best" approach, this distributes load
// across all healthy proxies while still favouring better-performing ones.
// This also avoids the thundering-herd problem where all requests pile onto
// a single proxy just because it has the highest score.
type WeightedSelector struct {
	randFloat64 func() float64
}

// NewWeightedRandom initializes a new selector with the standard
// math/rand generator.
func NewWeightedRandom() *WeightedSelector {
	return &WeightedSelector{randFloat64: rand.Float64}
}

// Select returns a randomly chosen entry weighted by Stats.Weight.
// Entries with weight 0 are effectively excluded from selection.
func (w *WeightedSelector) Select(entries []*Entry) *Entry {
	if len(entries) == 0 {
		return nil
	}

	var sum float64
	for _, entry := range entries {
		sum += entry.stats.Weight()
	}

	r := w.randFloat64() * sum
	for _, entry := range entries {
		r -= entry.stats.Weight()
		if r <= 0 {
			return entry
		}
	}

	return entries[len(entries)-1]
}
