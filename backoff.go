package client

import "time"

type BackoffStrategy interface {
	Next(attempt int) time.Duration
}
