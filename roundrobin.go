package client

import "sync/atomic"

type RoundRobin struct {
	counter atomic.Uint64
}

func (r *RoundRobin) Select(entries []*Entry) *Entry {
	if len(entries) == 0 {
		return nil
	}

	index := r.counter.Add(1) % uint64(len(entries))
	return entries[index]
}
