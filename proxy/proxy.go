package proxy

import (
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/gocolly/colly"
)

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (r *roundRobinSwitcher) GetProxy(_ *http.Request) (*url.URL, error) {
	u := r.proxyURLs[r.index%uint32(len(r.proxyURLs))]
	atomic.AddUint32(&r.index, 1)
	return u, nil
}

// RoundRobinProxySwitcher creates a proxy switcher function which rotates
// ProxyURLs on every request.
// The proxy type is determined by the URL scheme. "http"
// and "socks5" are supported. If the scheme is empty,
// "http" is assumed.
func RoundRobinProxySwitcher(ProxyURLs ...string) (colly.ProxyFunc, error) {
	urls := make([]*url.URL, len(ProxyURLs))
	for i, u := range ProxyURLs {
		parsedU, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsedU
	}
	return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}
