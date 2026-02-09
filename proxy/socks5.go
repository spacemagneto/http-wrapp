package proxy

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type SOCKS5Proxy struct {
	url string
}

func NewSOCKS5Proxy(proxyURL string) (*SOCKS5Proxy, error) {
	parseURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(parseURL.Scheme, "socks5") {
		return nil, fmt.Errorf("invalid proxy scheme: expected 'socks5://', got '%s://'", parseURL.Scheme)
	}

	return &SOCKS5Proxy{url: proxyURL}, nil
}

func (s *SOCKS5Proxy) Dial() fasthttp.DialFunc {
	return fasthttpproxy.FasthttpSocksDialer(s.url)
}
