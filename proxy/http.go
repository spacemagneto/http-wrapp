package proxy

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
)

type HTTPProxy struct {
	url     string
	timeout time.Duration
}

func NewHTTPProxy(proxyURL string, timeout time.Duration) (*HTTPProxy, error) {
	parseURL, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(parseURL.Scheme, "http") {
		return nil, fmt.Errorf("invalid proxy scheme: expected 'http://' or 'https://', got '%s://'", parseURL.Scheme)
	}

	return &HTTPProxy{url: proxyURL, timeout: timeout}, nil
}

func (h *HTTPProxy) Dial() fasthttp.DialFunc {
	if h.timeout > 0 {
		return fasthttpproxy.FasthttpHTTPDialerDualStackTimeout(h.url, h.timeout)
	}

	return fasthttpproxy.FasthttpHTTPDialer(h.url)
}
