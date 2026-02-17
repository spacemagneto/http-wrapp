package client

import "time"

type BackoffStrategy interface {
	Next(int64) time.Duration
}
