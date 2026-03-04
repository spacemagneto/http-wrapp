package client

type Entry struct {
	proxy Proxy
	stats Stats
}

func newEntry(proxy Proxy) *Entry {
	return &Entry{proxy: proxy}
}

func (e *Entry) Proxy() Proxy {
	return e.proxy
}

func (e *Entry) Stats() *Stats {
	return &e.stats
}
