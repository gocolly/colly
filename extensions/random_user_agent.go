package extensions

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/gocolly/colly/v2"
)

var uaGens = []func() string{
	genFirefoxUA,
	genChromeUA,
	genEdgeUA,
	genOperaUA,
}

var uaGensMobile = []func() string{
	genMobileUcwebUA,
	genMobileNexus10UA,
}

// RandomUserAgent generates a random DESKTOP browser user-agent on every requests
func RandomUserAgent(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", uaGens[rand.Intn(len(uaGens))]())
	})
}

// RandomMobileUserAgent generates a random MOBILE browser user-agent on every requests
func RandomMobileUserAgent(c *colly.Collector) {
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", uaGensMobile[rand.Intn(len(uaGensMobile))]())
	})
}

var ffVersions = []float32{
	// 2020
	72.0,
	73.0,
	74.0,
	75.0,
	76.0,
	77.0,
	78.0,
	79.0,
	80.0,
	81.0,
	82.0,
	83.0,
	84.0,

	// 2021
	85.0,
	86.0,
	87.0,
}

var chromeVersions = []string{
	// 2020
	"79.0.3945.117",
	"79.0.3945.130",
	"80.0.3987.106",
	"80.0.3987.116",
	"80.0.3987.122",
	"80.0.3987.132",
	"80.0.3987.149",
	"80.0.3987.163",
	"80.0.3987.87",
	"81.0.4044.113",
	"81.0.4044.122",
	"81.0.4044.129",
	"81.0.4044.138",
	"81.0.4044.92",
	"83.0.4103.106",
	"83.0.4103.116",
	"83.0.4103.97",
	"84.0.4147.105",
	"84.0.4147.125",
	"84.0.4147.135",
	"85.0.4183.102",
	"85.0.4183.121",
	"85.0.4183.83",
	"86.0.4240.111",
	"86.0.4240.183",
	"86.0.4240.198",
	"86.0.4240.75",

	// 2021
	"87.0.4280.141",
	"87.0.4280.66",
	"87.0.4280.88",
	"88.0.4324.146",
	"88.0.4324.182",
	"88.0.4324.190",
	"89.0.4389.114",
	"89.0.4389.90",
	"90.0.4430.72",
}

var edgeVersions = []string{
	"79.0.3945.74,79.0.309.43",
	"80.0.3987.87,80.0.361.48",
	"84.0.4147.105,84.0.522.50",
	"89.0.4389.128,89.0.774.77",
	"90.0.4430.72,90.0.818.39",
}

var operaVersions = []string{
	"2.7.62 Version/11.00",
	"2.2.15 Version/10.10",
	"2.9.168 Version/11.50",
	"2.2.15 Version/10.00",
	"2.8.131 Version/11.11",
	"2.5.24 Version/10.54",
}

var ucwebVersions = []string{
	"10.9.8.1006",
	"11.0.0.1016",
	"11.0.6.1040",
	"11.1.0.1041",
	"11.1.1.1091",
	"11.1.2.1113",
	"11.1.3.1128",
	"11.2.0.1125",
	"11.3.0.1130",
	"11.4.0.1180",
	"11.4.1.1138",
	"11.5.2.1188",
}

var androidVersions = []string{
	"4.4.2",
	"4.4.4",
	"5.0",
	"5.0.1",
	"5.0.2",
	"5.1",
	"5.1.1",
	"5.1.2",
	"6.0",
	"6.0.1",
	"7.0",
	"7.1.1",
	"7.1.2",
	"8.0.0",
	"8.1.0",
	"9",
	"10",
	"11",
}

var ucwebDevices = []string{
	"SM-C111",
	"SM-J727T1",
	"SM-J701F",
	"SM-J330G",
	"SM-N900",
	"DLI-TL20",
	"LG-X230",
	"AS-5433_Secret",
	"IdeaTabA1000-G",
	"GT-S5360",
	"HTC_Desire_601_dual_sim",
	"ALCATEL_ONE_TOUCH_7025D",
	"SM-N910H",
	"Micromax_Q4101",
	"SM-G600FY",
}

var nexus10Builds = []string{
	"JOP40D",
	"JOP40F",
	"JVP15I",
	"JVP15P",
	"JWR66Y",
	"KTU84P",
	"LMY47D",
	"LMY47V",
	"LMY48M",
	"LMY48T",
	"LMY48X",
	"LMY49F",
	"LMY49H",
	"LRX21P",
	"NOF27C",
}

var nexus10Safari = []string{
	"534.30",
	"535.19",
	"537.22",
	"537.31",
	"537.36",
	"600.1.4",
}

var osStrings = []string{
	// MacOS - High Sierra
	"Macintosh; Intel Mac OS X 10_13",
	"Macintosh; Intel Mac OS X 10_13_1",
	"Macintosh; Intel Mac OS X 10_13_2",
	"Macintosh; Intel Mac OS X 10_13_3",
	"Macintosh; Intel Mac OS X 10_13_4",
	"Macintosh; Intel Mac OS X 10_13_5",
	"Macintosh; Intel Mac OS X 10_13_6",

	// MacOS - Mojave
	"Macintosh; Intel Mac OS X 10_14",
	"Macintosh; Intel Mac OS X 10_14_1",
	"Macintosh; Intel Mac OS X 10_14_2",
	"Macintosh; Intel Mac OS X 10_14_3",
	"Macintosh; Intel Mac OS X 10_14_4",
	"Macintosh; Intel Mac OS X 10_14_5",
	"Macintosh; Intel Mac OS X 10_14_6",

	// MacOS - Catalina
	"Macintosh; Intel Mac OS X 10_15",
	"Macintosh; Intel Mac OS X 10_15_1",
	"Macintosh; Intel Mac OS X 10_15_2",
	"Macintosh; Intel Mac OS X 10_15_3",
	"Macintosh; Intel Mac OS X 10_15_4",
	"Macintosh; Intel Mac OS X 10_15_5",
	"Macintosh; Intel Mac OS X 10_15_6",
	"Macintosh; Intel Mac OS X 10_15_7",

	// MacOS - Big Sur
	"Macintosh; Intel Mac OS X 11_0",
	"Macintosh; Intel Mac OS X 11_0_1",
	"Macintosh; Intel Mac OS X 11_1",
	"Macintosh; Intel Mac OS X 11_2",
	"Macintosh; Intel Mac OS X 11_2_1",
	"Macintosh; Intel Mac OS X 11_2_2",
	"Macintosh; Intel Mac OS X 11_2_3",

	// Windows
	"Windows NT 10.0; Win64; x64",
	"Windows NT 5.1",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",

	// Linux
	"X11; Linux x86_64",
}

// Generates Firefox Browser User-Agent (Desktop)
//	-> "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:87.0) Gecko/20100101 Firefox/87.0"
func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

// Generates Chrome Browser User-Agent (Desktop)
//	-> "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"
func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

// Generates Microsoft Edge User-Agent (Desktop)
//	-> "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36 Edg/90.0.818.39"
func genEdgeUA() string {
	version := edgeVersions[rand.Intn(len(edgeVersions))]
	chrome_version := strings.Split(version, ",")[0]
	edge_version := strings.Split(version, ",")[1]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s", os, chrome_version, edge_version)
}

// Generates Opera Browser User-Agent (Desktop)
//	-> "Opera/9.80 (X11; Linux x86_64; U; en) Presto/2.8.131 Version/11.11"
func genOperaUA() string {
	version := operaVersions[rand.Intn(len(operaVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Opera/9.80 (%s; U; en) Presto/%s", os, version)
}

// Generates UCWEB/Nokia203 Browser User-Agent (Mobile)
//	-> "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCMini/10.9.8.1006 (SpeedMode; Proxy; Android 4.4.4; SM-J110H ) U2/1.0.0 Mobile"
func genMobileUcwebUA() string {
	device := ucwebDevices[rand.Intn(len(ucwebDevices))]
	version := ucwebVersions[rand.Intn(len(ucwebVersions))]
	android := androidVersions[rand.Intn(len(androidVersions))]
	return fmt.Sprintf("UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCMini/%s (SpeedMode; Proxy; Android %s; %s ) U2/1.0.0 Mobile", version, android, device)
}

// Generates Nexus 10 Browser User-Agent (Mobile)
//	-> "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 10 Build/LMY48T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.91 Safari/537.36"
func genMobileNexus10UA() string {
	build := nexus10Builds[rand.Intn(len(nexus10Builds))]
	android := androidVersions[rand.Intn(len(androidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	safari := nexus10Safari[rand.Intn(len(nexus10Safari))]
	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Nexus 10 Build/%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/%s", android, build, chrome, safari)
}
