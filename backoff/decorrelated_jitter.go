package backoff

import (
	"math/rand/v2"
	"time"
)

type DecorrelatedJitter struct {
	baseDelay, maxDelay time.Duration
}

func NewDecorrelatedJitter(delay, max time.Duration) *DecorrelatedJitter {
	if delay <= 0 {
		delay = DefaultJitterBaseDelay
	}

	if max <= 0 {
		max = DefaultJitterMaxDelay
	}

	if max < delay {
		max = delay
	}

	return &DecorrelatedJitter{baseDelay: delay, maxDelay: max}
}

func (d *DecorrelatedJitter) Next(previousDelay int64) time.Duration {
	// sleep = min(cap, random_between(base, sleep * 3))
	high := time.Duration(previousDelay) * 3
	if high < d.baseDelay {
		high = d.baseDelay
	}

	if high > d.maxDelay {
		high = d.maxDelay
	}

	diff := int64(high - d.baseDelay)
	if diff <= 0 {
		return d.baseDelay
	}

	return d.baseDelay + time.Duration(rand.Int64N(diff+1))
}
