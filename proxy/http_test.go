package proxy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
