package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestClient(t *testing.T) {
	t.Parallel()

	defaultClient := &fasthttp.Client{
		ReadTimeout:                   500 * time.Millisecond,
		WriteTimeout:                  500 * time.Millisecond,
		MaxIdleConnDuration:           1 * time.Hour,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	baseClient := NewClient(defaultClient)
	assert.NotNil(t, baseClient)

	t.Run("ServerError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		resp, err := baseClient.Get(ts.URL)
		assert.NoError(t, err)
		defer fasthttp.ReleaseResponse(resp)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode())
	})

	t.Run("InvalidURL", func(t *testing.T) {
		resp, err := baseClient.Get("http://localhost:1")
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
