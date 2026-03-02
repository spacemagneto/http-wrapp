package client

import (
	"sync/atomic"
	"time"
)

type Stats struct {
	// consecutiveFails counts failures in a row since the last success.
	// This is the value checked for quarantine: it resets on every success.
	consecutiveFails atomic.Int32

	failCount atomic.Int64

	successCount atomic.Int64

	// lastFailedUnix stores the UnixNano timestamp of the most recent failure.
	// Used by Entry.IsHealthy to determine if the cooldown window has elapsed.
	lastFailedUnix atomic.Int64
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) RecordSuccess() {
	s.successCount.Add(1)
	s.consecutiveFails.Store(0)
}

func (s *Stats) RecordFailed() {
	s.consecutiveFails.Add(1)
	s.failCount.Add(1)
	s.lastFailedUnix.Store(time.Now().UnixNano())
}
