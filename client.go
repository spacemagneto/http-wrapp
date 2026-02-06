package client

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type Client struct {
	client *fasthttp.Client
}

func NewClient(client *fasthttp.Client) *Client {
	return &Client{client: client}
}

func (c *Client) Get(url string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	res := fasthttp.AcquireResponse()

	if err := c.client.Do(req, res); err != nil {
		fasthttp.ReleaseResponse(res)
		return nil, fmt.Errorf("failed request: %w", err)
	}

	return res, nil
}
