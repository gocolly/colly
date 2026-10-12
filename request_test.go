package colly

import (
	"bytes"
	"net/url"
	"testing"
)

func TestNewRequest(t *testing.T) {

	ctx := &Context{}
	host := "test host"
	headerK := "headerk"
	headerV := "headerv"
	depth := 8
	method := "POST"
	body := bytes.NewReader([]byte{1, 2, 3})
	respCharSet := "testC"
	proxyUrl := "testUrl"
	url := &url.URL{Path: "test"}
	var id uint32 = 123

	req := NewRequest(
		url,
		WithRequestContext(ctx),
		WithRequestSetHeader(headerK, headerV),
		WithRequestDepth(depth),
		WithRequestMethod(method),
		WithRequestBody(body),
		WithRequestResponseCharacterEncoding(respCharSet),
		WithRequestProxyURL(proxyUrl),
		WithRequestHost(host),
		WithRequestID(id),
	)

	if req.URL != url || req.Ctx != ctx || req.Headers.Get(headerK) != headerV || req.Depth != depth || req.Method != req.Method || req.Body != body || req.ResponseCharacterEncoding != respCharSet || req.ProxyURL != proxyUrl || req.ID != id {
		t.Fail()
	}

}
