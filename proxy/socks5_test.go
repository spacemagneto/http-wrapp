package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/things-go/go-socks5"
	"github.com/valyala/fasthttp"
)

func TestSocks5Proxy(t *testing.T) {
	t.Parallel()

	t.Run("ValidWithoutAuth", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("socks5://127.0.0.1:1080")

		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		assert.Equal(t, "socks5://127.0.0.1:1080", proxy.url)
		assert.Empty(t, proxy.username)
		assert.Empty(t, proxy.password)
	})

	t.Run("ValidWithUsernameOnly", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("socks5://alice@127.0.0.1:1080")

		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		assert.Equal(t, "socks5://alice@127.0.0.1:1080", proxy.url)
		assert.Equal(t, "alice", proxy.username)
		assert.Empty(t, proxy.password)
	})

	t.Run("ValidWithUsernameAndPassword", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("socks5://alice:secret123@192.168.1.55:9050")

		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		assert.Equal(t, "alice", proxy.username)
		assert.Equal(t, "secret123", proxy.password)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("http://127.0.0.1:1080")

		assert.Error(t, err)
		assert.Nil(t, proxy)
	})

	t.Run("EmptyURL", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("")

		assert.Error(t, err)
		assert.Nil(t, proxy)
	})
}

func TestNewSOCKS5ProxyWithAuth(t *testing.T) {
	t.Parallel()

	t.Run("ValidURLWithAuth", func(t *testing.T) {
		proxy, err := NewSOCKS5ProxyWithAuth("socks5://127.0.0.1:1", "spacemagneto", "123456789")

		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		assert.Equal(t, "socks5://127.0.0.1:1", proxy.url)
		assert.Equal(t, "spacemagneto", proxy.username)
		assert.Equal(t, "123456789", proxy.password)
	})

	t.Run("AuthWithOverridesAuthData", func(t *testing.T) {
		proxy, err := NewSOCKS5ProxyWithAuth("socks5://spacemagneto:123456789@127.0.0.1:1", "spacemagneto_second", "1234567890")

		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		assert.Equal(t, "spacemagneto_second", proxy.username)
		assert.Equal(t, "1234567890", proxy.password)
	})

	t.Run("InvalidURL", func(t *testing.T) {
		proxy, err := NewSOCKS5ProxyWithAuth("http://127.0.0.1:1", "user", "pass")

		assert.Error(t, err)
		assert.Nil(t, proxy)
	})
}

func TestSOCKS5ProxyDial(t *testing.T) {
	t.Parallel()

	t.Run("SuccessInitDialFunc", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("socks5://127.0.0.1:1080")
		assert.NoError(t, err)

		assert.NotNil(t, proxy.Dial())
	})

	t.Run("FailedInitDialFuncWithInvalidProxyURL", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("socks5://127.0.0.1:1")
		assert.NoError(t, err)

		dialFunc := proxy.Dial()

		conn, err := dialFunc("google.com")

		assert.Error(t, err)
		assert.Nil(t, conn)
	})

	t.Run("SuccessInitDialFuncWithFastHTTPClient", func(t *testing.T) {
		proxy, err := NewSOCKS5Proxy("socks5://user:pass@127.0.0.1:1080")
		assert.NoError(t, err)

		client := &fasthttp.Client{ReadTimeout: 300 * time.Millisecond, WriteTimeout: 300 * time.Millisecond, Dial: proxy.Dial()}
		assert.NotNil(t, client.Dial)

		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		req.SetRequestURI("https://google.com")

		err = client.Do(req, res)
		assert.Error(t, err)
	})
}

func TestSOCKS5ProxyDialWithSocksServer(t *testing.T) {
	socksServer := socks5.NewServer()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)

	defer listener.Close()

	go socksServer.Serve(listener)

	time.Sleep(100 * time.Millisecond)

	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))
	defer targetServer.Close()

	proxy, err := NewSOCKS5Proxy(fmt.Sprintf("socks5://%s", listener.Addr().String()))
	assert.NoError(t, err)

	client := &fasthttp.Client{ReadTimeout: 300 * time.Millisecond, WriteTimeout: 300 * time.Millisecond, Dial: proxy.Dial()}
	assert.NotNil(t, proxy.Dial)

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(targetServer.URL)

	err = client.Do(req, res)

	assert.NoError(t, err)
	assert.Equal(t, fasthttp.StatusOK, res.StatusCode())
	assert.Equal(t, "success", string(res.Body()))
}
