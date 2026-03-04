package client

import "time"

type PoolConfig struct {
	MaxFails int64
	Cooldown time.Duration
	Selector Selector
}

func defaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxFails: 0,
		Cooldown: 0,
		Selector: nil,
	}
}

type Pool struct {
}

func NewPool() *Pool {
	return &Pool{}
}
