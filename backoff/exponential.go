package backoff

import (
	"math"
	"time"
)

type Exponential struct {
	baseDelay, maxDelay time.Duration
	step                float64
}

func NewExponential(delay, max time.Duration) *Exponential {
	return &Exponential{baseDelay: delay, maxDelay: max, step: 2}
}

func NewExponentialWithStep(delay, max time.Duration, step float64) *Exponential {
	if step <= 0 {
		step = 2
	}

	return &Exponential{baseDelay: delay, maxDelay: max, step: step}
}

func (e *Exponential) Next(attempt int) time.Duration {
	// calculate the delay: base * step^attempt
	delay := time.Duration(float64(e.baseDelay) * math.Pow(e.step, float64(attempt)))

	if delay > e.maxDelay || delay < 0 {
		return e.maxDelay
	}

	return delay
}
