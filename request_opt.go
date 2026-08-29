package colly

import (
	"io"
	"net/http"
)

type RequestOpt func(req *Request)

func WithRequestProxyURL(proxyUrl string) RequestOpt {
	return func(req *Request) {
		req.ProxyURL = proxyUrl
	}
}
func WithRequestID(id uint32) RequestOpt {
	return func(req *Request) {
		req.ID = id
	}
}
func WithRequestResponseCharacterEncoding(responseCharacterEncoding string) RequestOpt {
	return func(req *Request) {
		req.ResponseCharacterEncoding = responseCharacterEncoding
	}
}

func WithRequestBody(body io.Reader) RequestOpt {
	return func(req *Request) {
		req.Body = body
	}
}

func WithRequestMethod(method string) RequestOpt {
	return func(req *Request) {
		req.Method = method
	}
}
func WithRequestHost(host string) RequestOpt {
	return func(req *Request) {
		req.Host = host
	}
}

func WithRequestDepth(depth int) RequestOpt {
	return func(req *Request) {
		req.Depth = depth
	}
}

func WithRequestSetHeader(k, v string) RequestOpt {
	return func(req *Request) {
		if req.Headers == nil {
			req.Headers = &http.Header{}
		}
		req.Headers.Add(k, v)
	}
}
func WithRequestContext(ctx *Context) RequestOpt {
	return func(req *Request) {
		req.Ctx = ctx
	}
}
