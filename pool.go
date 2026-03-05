package client

import "time"

type PoolConfig struct {
	MaxFails       int64
	CooldownWindow time.Duration
	Selector       Selector
}

func defaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxFails:       0,
		CooldownWindow: 0,
		Selector:       nil,
	}
}

type Pool struct {
	entries []*Entry
}

func NewPool(proxies []Proxy, cfg *PoolConfig) *Pool {
	defaultCfg := defaultPoolConfig()

	if cfg.MaxFails == 0 {
		cfg.MaxFails = defaultCfg.MaxFails
	}

	if cfg.CooldownWindow == 0 {
		cfg.CooldownWindow = defaultCfg.CooldownWindow
	}

	if cfg.Selector == nil {
		cfg.Selector = defaultCfg.Selector
	}

	entries := make([]*Entry, len(proxies))
	for _, proxy := range proxies {
		entries = append(entries, newEntry(proxy))
	}

	return &Pool{
		entries: entries,
	}
}
