package client

import (
	"sync"
	"time"
)

type PoolConfig struct {
	MaxFails       int64
	CooldownWindow time.Duration
	Selector       Selector
}

func defaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxFails:       3,
		CooldownWindow: 30 * time.Second,
		Selector:       &RoundRobinSelector{},
	}
}

type Pool struct {
	mutex   sync.RWMutex
	entries []*Entry
	cfg     PoolConfig
}

func NewPool(proxies []Proxy, cfg PoolConfig) *Pool {
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

	entries := make([]*Entry, 0, len(proxies))
	for _, proxy := range proxies {
		entries = append(entries, newEntry(proxy))
	}

	return &Pool{entries: entries, cfg: cfg}
}

func (p *Pool) Pick() (*Entry, error) {
	p.mutex.RLock()
	healthy := p.healthyEntries()
	all := p.entries
	p.mutex.RUnlock()

	entriesList := healthy
	if len(entriesList) == 0 {
		if len(all) == 0 {
			return nil, ErrNoProxies
		}

		entriesList = all
	}

	entry := p.cfg.Selector.Select(entriesList)
	return entry, nil
}

func (p *Pool) healthyEntries() []*Entry {
	healthyProxies := make([]*Entry, 0)

	for _, entry := range p.entries {
		if entry.HealthCheck(p.cfg.MaxFails, p.cfg.CooldownWindow) {
			healthyProxies = append(healthyProxies, entry)
		}
	}

	return healthyProxies
}
