package client

import "sync/atomic"

type Stats struct {
	// consecutiveFails counts failures in a row since the last success.
	// This is the value checked for quarantine: it resets on every success.
	consecutiveFails atomic.Int32

	failsCount atomic.Int64

	successCount atomic.Int64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) RecordSuccess() {
	s.successCount.Add(1)
	s.consecutiveFails.Store(0)
}
