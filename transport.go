package client

import "github.com/valyala/fasthttp"

type Transport interface {
	Do(req *fasthttp.Request, res *fasthttp.Response) error
}
