package backoff

import (
	"math/rand"
	"time"
)

const (
	DefaultJitterBaseDelay = 200 * time.Millisecond
	DefaultJitterMaxDelay  = 30 * time.Second
)

type Jitter struct {
	baseDelay, maxDelay time.Duration
}

func NewJitter(delay, max time.Duration) *Jitter {
	if delay <= 0 {
		delay = DefaultJitterBaseDelay
	}

	if max <= 0 {
		max = DefaultJitterMaxDelay
	}

	if max < delay {
		max = delay
	}

	return &Jitter{baseDelay: delay, maxDelay: max}
}

func (j *Jitter) Next(attempt int) time.Duration {
	// sleep = random_between(0, min(cap, base * 2 ** attempt))
	// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
	limit := int64(j.baseDelay) * (1 << uint64(attempt))

	if limit > int64(j.maxDelay) || limit <= 0 {
		limit = int64(j.maxDelay)
	}

	return time.Duration(rand.Int63n(limit))
}
