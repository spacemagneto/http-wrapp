package client

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
