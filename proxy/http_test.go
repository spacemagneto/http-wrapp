package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestHTTPProxy(t *testing.T) {
	t.Parallel()

	t.Run("ValidWithoutAuth", func(t *testing.T) {
		proxy, err := NewHTTPProxy("https://127.0.0.1:1080", 3*time.Second)

		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		assert.Equal(t, "https://127.0.0.1:1080", proxy.url)
		assert.Equal(t, 3*time.Second, proxy.timeout)
	})

	t.Run("ValidWithUsernameOnly", func(t *testing.T) {
		proxy, err := NewHTTPProxy("https://alice@127.0.0.1:1080", 0)

		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		assert.Equal(t, "https://alice@127.0.0.1:1080", proxy.url)
		assert.Equal(t, time.Duration(0), proxy.timeout)
	})

	t.Run("ValidWithUsernameAndPassword", func(t *testing.T) {
		proxy, err := NewHTTPProxy("https://alice:secret123@192.168.1.55:9050", 0)

		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		assert.Equal(t, "https://alice:secret123@192.168.1.55:9050", proxy.url)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("socks4://127.0.0.1:1080")

		assert.Error(t, err)
		assert.Nil(t, proxy)
		assert.Contains(t, err.Error(), "invalid proxy scheme")
	})

	t.Run("EmptyURL", func(t *testing.T) {
		proxy, err := NewHTTPProxy("", 0)

		assert.Error(t, err)
		assert.Nil(t, proxy)
	})
}

func TestHTTPProxyDial(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitDialFunc", func(t *testing.T) {
		targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("hello from target"))
		}))
		defer targetServer.Close()

		proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test-Proxy", "http_proxy")

			resp, _ := http.Get(targetServer.URL)
			w.WriteHeader(resp.StatusCode)
			_ = resp.Write(w)
		}))
		defer proxyServer.Close()

		proxy, err := NewHTTPProxy(proxyServer.URL, 2*time.Second)
		assert.NoError(t, err)

		client := &fasthttp.Client{ReadTimeout: 300 * time.Millisecond, WriteTimeout: 300 * time.Millisecond, Dial: proxy.Dial()}
		assert.NotNil(t, client.Dial)

		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		req.SetRequestURI(targetServer.URL)
		err = client.Do(req, res)

		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode())
		assert.Equal(t, "http_proxy", string(res.Header.Peek("X-Test-Proxy")))
	})
}
