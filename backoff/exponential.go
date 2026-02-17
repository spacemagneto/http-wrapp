package backoff

import (
	"math"
	"time"
)

const (
	DefaultBaseDelay = 200 * time.Millisecond
	DefaultMaxDelay  = 10 * time.Second
	DefaultStep      = 2.0
)

type Exponential struct {
	baseDelay, maxDelay time.Duration
	step                float64
}

func NewExponential(delay, max time.Duration) *Exponential {
	return NewExponentialWithStep(delay, max, DefaultStep)
}

func NewExponentialWithStep(delay, max time.Duration, step float64) *Exponential {
	if delay <= 0 {
		delay = DefaultBaseDelay
	}

	if max <= 0 {
		max = DefaultMaxDelay
	}

	if max < delay {
		max = delay
	}

	if step <= 1.0 {
		step = DefaultStep
	}

	return &Exponential{baseDelay: delay, maxDelay: max, step: step}
}

func (e *Exponential) Next(attempt int64) time.Duration {
	// calculate the delay: base * step^attempt
	delay := time.Duration(float64(e.baseDelay) * math.Pow(e.step, float64(attempt)))

	if delay > e.maxDelay || delay < 0 {
		return e.maxDelay
	}

	return delay
}
