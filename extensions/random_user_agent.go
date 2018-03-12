package extensions

import (
	"fmt"
	"math/rand"

	"github.com/gocolly/colly"
)

var uaGens = []func() string{
	genFirefoxUA,
}

// RandomUserAgent generates a random browser user agent on every request
func RandomUserAgent(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", uaGens[rand.Intn(len(uaGens))]())
	})
}

var ffVersions = []float32{
	40.0,
	39.0,
	38.0,
	37.0,
	36.0,
	35.0,
}

var ffOSs = []string{
	"Macintosh; Intel Mac OS X 10_10",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",
	"X11; Linux x86_64",
}

func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := ffOSs[rand.Intn(len(ffOSs))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}
