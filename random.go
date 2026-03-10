package client

import "math/rand/v2"

// RandomSelector implements a uniform selection strategy where each entry has an equal
// probability of being chosen.
//
// Use this strategy when performance metrics (like latency or success rate)
// are either unavailable or when the underlying resources are expected
// to perform identically. It provides the lowest possible computational
// overhead for a selection process.
type RandomSelector struct{}

// Select picks an entry from the slice with uniform probability.
// This operation is O(1) and does not account for any external statistics.
func (r *RandomSelector) Select(entries []*Entry) *Entry {
	if len(entries) == 0 {
		return nil
	}

	return entries[rand.IntN(len(entries))]
}
