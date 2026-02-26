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

		nextEntry := roundRobin.Next(nil)
		assert.Nil(t, nextEntry)
		assert.Equal(t, uint64(0), roundRobin.counter.Load())
	})

	t.Run("SuccessGetNextProxy", func(t *testing.T) {
		roundRobin := &RoundRobin{}
		assert.NotNil(t, roundRobin)
		assert.Equal(t, uint64(0), roundRobin.counter.Load())

		entries := []*Entry{newEntry(&mockProxy{id: 1}), newEntry(&mockProxy{id: 2}), newEntry(&mockProxy{id: 3})}
		assert.Equal(t, entries[1].Proxy(), roundRobin.Next(entries).Proxy())
		assert.Equal(t, entries[2].Proxy(), roundRobin.Next(entries).Proxy())
		assert.Equal(t, entries[0].Proxy(), roundRobin.Next(entries).Proxy())
		assert.Equal(t, entries[1].Proxy(), roundRobin.Next(entries).Proxy())
	})
}
