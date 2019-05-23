package extensions

import (
	"fmt"
	"math/rand"

	"github.com/gocolly/colly"
)

var uaGens = []func() string{
	genFirefoxUA,
	genChromeUA,
	genOperaUA,
}

// RandomUserAgent generates a random browser user agent on every request
func RandomUserAgent(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", uaGens[rand.Intn(len(uaGens))]())
	})
}

var ffVersions = []float32{
	58.0,
	57.0,
	56.0,
	52.0,
	48.0,
	40.0,
	35.0,
}

var chromeVersions = []string{
	"65.0.3325.146",
	"64.0.3282.0",
	"41.0.2228.0",
	"40.0.2214.93",
	"37.0.2062.124",
}

var operaVersions = []string{
	"2.7.62 Version/11.00",
	"2.2.15 Version/10.10",
	"2.9.168 Version/11.50",
	"2.2.15 Version/10.00",
	"2.8.131 Version/11.11",
	"2.5.24 Version/10.54",
}

var osStrings = []string{
	"Macintosh; Intel Mac OS X 10_10",
	"Windows NT 10.0",
	"Windows NT 5.1",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",
	"X11; Linux x86_64",
}

func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

func genOperaUA() string {
	version := operaVersions[rand.Intn(len(operaVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Opera/9.80 (%s; U; en) Presto/%s", os, version)
}
