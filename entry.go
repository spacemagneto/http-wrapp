package client

type entry struct {
	proxy Proxy
	stats proxyStats
}

func newEntry(proxy Proxy) *entry {
	return &entry{proxy: proxy}
}

func (e *entry) Proxy() Proxy {
	return e.proxy
}
