package proxy

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type SOCKS5Proxy struct {
	url, username, password string
}

func NewSOCKS5Proxy(proxyURL string) (*SOCKS5Proxy, error) {
	usr, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(usr.Scheme, "socks5") {
		return nil, fmt.Errorf("invalid proxy scheme: expected 'socks5://', got '%s://'", usr.Scheme)
	}

	proxy := &SOCKS5Proxy{url: proxyURL}

	if usr.User != nil {
		proxy.username = usr.User.Username()
		proxy.password, _ = usr.User.Password()
	}

	return proxy, err
}

func NewSOCKS5ProxyWithAuth(proxyURL, username, password string) (*SOCKS5Proxy, error) {
	proxy, err := NewSOCKS5Proxy(proxyURL)
	if err != nil {
		return nil, err
	}

	proxy.username = username
	proxy.password = password

	return proxy, nil
}

func (s *SOCKS5Proxy) Dial() fasthttp.DialFunc {
	return fasthttpproxy.FasthttpSocksDialer(s.url)
}
