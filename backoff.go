package client

import "time"

type Backoff interface {
	Next(args int64) time.Duration
}
