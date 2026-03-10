package client

import "time"

// Entry is the internal unit of the pool.
// It binds a Proxy to its Stats and exposes health-check logic.
//
// Entry is the only place in the package that knows about both Proxy and Stats,
// keeping the two concerns separate while providing a convenient handle for the pool.
type Entry struct {
	proxy Proxy
	stats Stats
}

func newEntry(proxy Proxy) *Entry {
	return &Entry{proxy: proxy, stats: Stats{}}
}

func (e *Entry) Proxy() Proxy {
	return e.proxy
}

func (e *Entry) Stats() *Stats {
	return &e.stats
}

// HealthCheck reports whether this proxy is eligible to receive requests.
//
// A proxy is considered unhealthy when it has accumulated maxConsecutiveFails
// failures in a row and the cooldown window has not yet elapsed since the
// last failure. Once the cooldown expires the proxy is given a second chance
// automatically — no manual reset is required.
func (e *Entry) HealthCheck(maxFails int64, cooldown time.Duration) bool {
	if e.stats.ConsecutiveFails() < maxFails {
		return true
	}

	return time.Since(e.stats.LastFailedTime()) >= cooldown
}
