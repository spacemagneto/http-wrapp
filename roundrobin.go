package client

import "sync/atomic"

// RoundRobinSelector cycles through entries in order, wrapping around at the end.
// It distributes load evenly regardless of proxy performance.
//
// The internal counter is updated atomically, making this selector
// safe for concurrent use without a mutex.
type RoundRobinSelector struct {
	counter atomic.Uint64
}

// Select returns the next entry in round-robin order.
// Concurrent calls are safe; each caller receives a distinct index.
func (r *RoundRobinSelector) Select(entries []*Entry) *Entry {
	if len(entries) == 0 {
		return nil
	}

	index := r.counter.Add(1) % uint64(len(entries))
	return entries[index]
}
