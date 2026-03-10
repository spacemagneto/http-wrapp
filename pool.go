package client

import (
	"sync"
	"time"
)

// PoolConfig holds the tuning parameters for a Pool.
type PoolConfig struct {
	// MaxFails is the number of consecutive failures after which
	// a proxy is placed in quarantine and skipped by Pick.
	// Defaults to 3 if zero.
	MaxFails int64

	// CooldownWindow is the duration a proxy stays in quarantine before
	// being given another chance. Defaults to 30s if zero.
	CooldownWindow time.Duration

	// Selector determines which healthy proxy Pick should hand out.
	// Defaults to RoundRobinSelector if nil.
	Selector Selector
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

// NewPool creates a Pool from the provided proxies and config.
// Any zero-value field in cfg is replaced with its default.
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

// Pick selects the next proxy to use according to the configured Selector.
//
// Only healthy entries (those that pass HealthCheck) are offered to the Selector.
// If every proxy is currently in quarantine, Pick falls back to selecting from
// the full list rather than returning an error — this prevents a total stall
// when all proxies are temporarily degraded.
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
