package backoff

import (
	"time"
)

type Fixed struct {
	delay time.Duration
}

func NewFixed(delay time.Duration) *Fixed {
	if delay <= 0 {
		delay = 1 * time.Second
	}

	return &Fixed{delay: delay}
}

func (f *Fixed) Next(_ int) time.Duration {
	return f.delay
}
