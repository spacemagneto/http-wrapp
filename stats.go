package client

import "sync/atomic"

type Stats struct {
	failsCount atomic.Int64

	successCount atomic.Int64
}

func NewStats() *Stats {
	return &Stats{}
}
