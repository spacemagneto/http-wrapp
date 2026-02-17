package backoff

import (
	"math/rand"
	"time"
)

type EqualJitter struct {
	baseDelay, maxDelay time.Duration
}

func NewEqualJitter(delay, max time.Duration) *EqualJitter {
	if delay <= 0 {
		delay = DefaultJitterBaseDelay
	}

	if max <= 0 {
		max = DefaultJitterMaxDelay
	}

	if max < delay {
		max = delay
	}

	return &EqualJitter{baseDelay: delay, maxDelay: max}
}

func (e *EqualJitter) Next(attempt int64) time.Duration {
	// temp = min(cap, base * 2 ** attempt)
	// sleep = temp / 2 + random_between(0, temp / 2)
	// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
	limit := int64(e.baseDelay) * (1 << uint64(attempt))

	if limit > int64(e.maxDelay) || limit <= 0 {
		limit = int64(e.maxDelay)
	}

	temp := limit / 2

	var delay int64
	if temp > 0 {
		delay = rand.Int63n(temp)
	}

	return time.Duration(temp + delay)
}
