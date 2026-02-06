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
		if err != nil {
			t.Fatalf("Get() should not return error on 500, but got: %v", err)
		}
		defer fasthttp.ReleaseResponse(resp)

		if resp.StatusCode() != 500 {
			t.Errorf("expected 500, got %d", resp.StatusCode())
		}
	})
}
