package client

import (
	"sync/atomic"
	"time"
)

const baseWeight float64 = 1.0

// Stats tracks runtime metrics for a single proxy instance.
// It is designed to be embedded in Entry and must not be copied after first use.
//
// All fields are updated atomically, so Stats is safe for concurrent use
// without additional locking.
//
// Stats deliberately knows nothing about Proxy — it is a pure metrics store.
// The association between a proxy and its stats is the responsibility of Entry.
type Stats struct {
	// consecutiveFails counts failures in a row since the last success.
	// This is the value checked for quarantine: it resets on every success.
	consecutiveFails atomic.Int64

	// failCount is a monotonically increasing failure counter.
	// It is never reset and provides a full history of failures.
	failCount atomic.Int64

	successCount atomic.Int64

	// totalLatencyMs accumulates response times for average calculation.
	totalLatencyMs atomic.Int64

	// latencyCount is incremented each time RecordLatency is called.
	// Kept separate from successCount so latency can be recorded independently.
	latencyCount atomic.Int64

	// lastFailedUnix stores the UnixNano timestamp of the most recent failure.
	// Used by Entry.HealthCheck to determine if the cooldown window has elapsed.
	lastFailedUnix atomic.Int64
}

// RecordSuccess increments the success counter and resets consecutiveFails.
// Call this after every request that completes without a network-level error
// and returns a non-retryable HTTP status.
func (s *Stats) RecordSuccess() {
	s.successCount.Add(1)
	s.consecutiveFails.Store(0)
}

// RecordFailed increments both consecutiveFails and failCount,
// and timestamps the event. Call this on network errors, timeouts,
// and retryable HTTP responses (5xx, 429).
func (s *Stats) RecordFailed() {
	s.consecutiveFails.Add(1)
	s.failCount.Add(1)
	s.lastFailedUnix.Store(time.Now().UnixNano())
}

// RecordLatency adds a latency sample in milliseconds.
// Should be called alongside RecordSuccess to keep the average meaningful.
func (s *Stats) RecordLatency(ms int64) {
	s.totalLatencyMs.Add(ms)
	s.latencyCount.Add(1)
}

// ConsecutiveFails returns the number of failures since the last success.
// This is the primary signal used by HealthCheck to decide quarantine.
func (s *Stats) ConsecutiveFails() int64 {
	return s.consecutiveFails.Load()
}

// SuccessCount returns the total number of successful requests ever recorded.
func (s *Stats) SuccessCount() int64 {
	return s.successCount.Load()
}

// Failures returns the total number of failures ever recorded.
// Unlike ConsecutiveFails, this value never resets.
func (s *Stats) Failures() int64 {
	return s.failCount.Load()
}

// AvgLatencyMs returns the mean response time across all recorded samples.
// Returns 0 if no latency samples have been recorded yet.
func (s *Stats) AvgLatencyMs() float64 {
	count := s.latencyCount.Load()

	if count == 0 {
		return 0
	}

	return float64(s.totalLatencyMs.Load()) / float64(count)
}

// successRate returns the fraction of successful requests in range [0.0, 1.0].
// A proxy with no requests yet returns 1.0 (optimistic default / credit of trust),
// so new proxies are not penalised before they have had a chance to be used.
func (s *Stats) successRate() float64 {
	success := s.successCount.Load()
	total := success + s.failCount.Load()

	if total == 0 {
		return baseWeight
	}

	return float64(success) / float64(total)
}

// LastFailedTime returns the time of the most recent failure.
// Returns zero time if the proxy has never failed.
func (s *Stats) LastFailedTime() time.Time {
	ns := s.lastFailedUnix.Load()
	if ns == 0 {
		return time.Time{}
	}

	return time.Unix(0, ns)
}

// Weight computes a selection score used by WeightedSelector.
// The formula rewards both high success rate and low latency:
//
//	weight = successRate * (1000 / avgLatencyMs)
//
// A proxy with no latency data yet falls back to its success rate alone,
// which is 1.0 for brand-new proxies — giving them a fair initial chance.
func (s *Stats) Weight() float64 {
	rate := s.successRate()
	latency := s.AvgLatencyMs()

	if latency > 0 {
		return rate * (1000.0 / latency)
	}

	return rate
}
