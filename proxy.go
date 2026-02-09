package client

import "github.com/valyala/fasthttp"

type Proxy interface {
	Dial() fasthttp.DialFunc
}
