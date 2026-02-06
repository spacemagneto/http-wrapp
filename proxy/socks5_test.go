package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
