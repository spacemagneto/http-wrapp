package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

type mockProxy struct {
	id int
}

func (m *mockProxy) Dial() fasthttp.DialFunc {
	return nil
}

func TestRoundRobin(t *testing.T) {
	t.Parallel()

	t.Run("EmptyEntries", func(t *testing.T) {
		roundRobin := &RoundRobin{}
		assert.NotNil(t, roundRobin)
		assert.Equal(t, uint64(0), roundRobin.counter.Load())

		nextProxy := roundRobin.Next(nil)
		assert.Nil(t, nextProxy)
		assert.Equal(t, uint64(0), roundRobin.counter.Load())
	})

	t.Run("SuccessGetNextProxy", func(t *testing.T) {
		roundRobin := &RoundRobin{}
		assert.NotNil(t, roundRobin)
		assert.Equal(t, uint64(0), roundRobin.counter.Load())

		entries := []Proxy{&mockProxy{id: 1}, &mockProxy{id: 2}, &mockProxy{id: 3}}
		assert.Equal(t, entries[1], roundRobin.Next(entries))
		assert.Equal(t, entries[2], roundRobin.Next(entries))
		assert.Equal(t, entries[0], roundRobin.Next(entries))
		assert.Equal(t, entries[1], roundRobin.Next(entries))
	})
}
