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
