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
	// NOTE: Only version released after Jun 1, 2022 will be listed.
	// Data source: https://en.wikipedia.org/wiki/Firefox_version_history

	// 2022
	102.0,
	103.0,
	104.0,
	105.0,
	106.0,
	107.0,
	108.0,

	// 2023
	109.0,
	110.0,
	111.0,
	112.0,
}

var chromeVersions = []string{
	// NOTE: Only version released after Jun 1, 2022 will be listed.
	// Data source: https://chromereleases.googleblog.com/search/label/Stable%20updates

	// https://chromereleases.googleblog.com/2022/06/stable-channel-update-for-desktop.html
	"102.0.5005.115",

	// https://chromereleases.googleblog.com/2022/06/stable-channel-update-for-desktop_21.html
	"103.0.5060.53",

	// https://chromereleases.googleblog.com/2022/06/stable-channel-update-for-desktop_27.html
	"103.0.5060.66",

	// https://chromereleases.googleblog.com/2022/07/stable-channel-update-for-desktop.html
	"103.0.5060.114",

	// https://chromereleases.googleblog.com/2022/07/stable-channel-update-for-desktop_19.html
	"103.0.5060.134",

	// https://chromereleases.googleblog.com/2022/08/stable-channel-update-for-desktop.html
	"104.0.5112.79",
	"104.0.5112.80",
	"104.0.5112.81",

	// https://chromereleases.googleblog.com/2022/08/stable-channel-update-for-desktop_16.html
	"104.0.5112.101",
	"104.0.5112.102",

	// https://chromereleases.googleblog.com/2022/08/stable-channel-update-for-desktop_30.html
	"105.0.5195.52",
	"105.0.5195.53",
	"105.0.5195.54",

	// https://chromereleases.googleblog.com/2022/09/stable-channel-update-for-desktop.html
	"105.0.5195.102",

	// https://chromereleases.googleblog.com/2022/09/stable-channel-update-for-desktop_14.html
	"105.0.5195.125",
	"105.0.5195.126",
	"105.0.5195.127",

	// https://chromereleases.googleblog.com/2022/09/stable-channel-update-for-desktop_27.html
	"106.0.5249.61",
	"106.0.5249.62",

	// https://chromereleases.googleblog.com/2022/09/stable-channel-update-for-desktop_30.html
	"106.0.5249.91",

	// https://chromereleases.googleblog.com/2022/10/stable-channel-update-for-desktop.html
	"106.0.5249.103",

	// https://chromereleases.googleblog.com/2022/10/stable-channel-update-for-desktop_11.html
	"106.0.5249.119",

	// https://chromereleases.googleblog.com/2022/10/stable-channel-update-for-desktop_25.html
	"107.0.5304.62",
	"107.0.5304.63",
	"107.0.5304.68",

	// https://chromereleases.googleblog.com/2022/10/stable-channel-update-for-desktop_27.html
	"107.0.5304.87",
	"107.0.5304.88",

	// https://chromereleases.googleblog.com/2022/11/stable-channel-update-for-desktop.html
	"107.0.5304.106",
	"107.0.5304.107",
	"107.0.5304.110",

	// https://chromereleases.googleblog.com/2022/11/stable-channel-update-for-desktop_24.html
	"107.0.5304.121",
	"107.0.5304.122",

	// https://chromereleases.googleblog.com/2022/11/stable-channel-update-for-desktop_29.html
	"108.0.5359.71",
	"108.0.5359.72",

	// https://chromereleases.googleblog.com/2022/12/stable-channel-update-for-desktop.html
	"108.0.5359.94",
	"108.0.5359.95",

	// https://chromereleases.googleblog.com/2022/12/stable-channel-update-for-desktop_7.html
	"108.0.5359.98",
	"108.0.5359.99",

	// https://chromereleases.googleblog.com/2022/12/stable-channel-update-for-desktop_13.html
	"108.0.5359.124",
	"108.0.5359.125",

	// https://chromereleases.googleblog.com/2023/01/stable-channel-update-for-desktop.html
	"109.0.5414.74",
	"109.0.5414.75",
	"109.0.5414.87",

	// https://chromereleases.googleblog.com/2023/01/stable-channel-update-for-desktop_24.html
	"109.0.5414.119",
	"109.0.5414.120",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-update-for-desktop.html
	"110.0.5481.77",
	"110.0.5481.78",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update.html
	"110.0.5481.96",
	"110.0.5481.97",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update_14.html
	"110.0.5481.100",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update_16.html
	"110.0.5481.104",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update_22.html
	"110.0.5481.177",
	"110.0.5481.178",

	// https://chromereleases.googleblog.com/2023/02/stable-channel-desktop-update_97.html
	"109.0.5414.129",

	// https://chromereleases.googleblog.com/2023/03/stable-channel-update-for-desktop.html
	"111.0.5563.64",
	"111.0.5563.65",

	// https://chromereleases.googleblog.com/2023/03/stable-channel-update-for-desktop_21.html
	"111.0.5563.110",
	"111.0.5563.111",

	// https://chromereleases.googleblog.com/2023/03/stable-channel-update-for-desktop_27.html
	"111.0.5563.146",
	"111.0.5563.147",

	// https://chromereleases.googleblog.com/2023/04/stable-channel-update-for-desktop.html
	"112.0.5615.49",
	"112.0.5615.50",

	// https://chromereleases.googleblog.com/2023/04/stable-channel-update-for-desktop_12.html
	"112.0.5615.86",
	"112.0.5615.87",

	// https://chromereleases.googleblog.com/2023/04/stable-channel-update-for-desktop_14.html
	"112.0.5615.121",

	// https://chromereleases.googleblog.com/2023/04/stable-channel-update-for-desktop_18.html
	"112.0.5615.137",
	"112.0.5615.138",
	"112.0.5615.165",
}

var edgeVersions = []string{
	// NOTE: Only version released after Jun 1, 2022 will be listed.
	// Data source: https://learn.microsoft.com/en-us/deployedge/microsoft-edge-release-schedule

	// 2022
	"103.0.0.0,103.0.1264.37",
	"104.0.0.0,104.0.1293.47",
	"105.0.0.0,105.0.1343.25",
	"106.0.0.0,106.0.1370.34",
	"107.0.0.0,107.0.1418.24",
	"108.0.0.0,108.0.1462.42",

	// 2023
	"109.0.0.0,109.0.1518.49",
	"110.0.0.0,110.0.1587.41",
	"111.0.0.0,111.0.1661.41",
	"112.0.0.0,112.0.1722.34",
	"112.0.0.0,112.0.1722.54",
	"112.0.0.0,112.0.1722.58",
}

var operaVersions = []string{
	// NOTE: Only version released after Jan 1, 2023 will be listed.
	// Data source: https://blogs.opera.com/desktop/

	// https://blogs.opera.com/desktop/changelog-for-96/
	"110.0.5449.0,96.0.4640.0",
	"110.0.5464.2,96.0.4653.0",
	"110.0.5464.2,96.0.4660.0",
	"110.0.5481.30,96.0.4674.0",
	"110.0.5481.30,96.0.4691.0",
	"110.0.5481.30,96.0.4693.12",
	"110.0.5481.77,96.0.4693.16",
	"110.0.5481.100,96.0.4693.20",
	"110.0.5481.178,96.0.4693.31",
	"110.0.5481.178,96.0.4693.50",
	"110.0.5481.192,96.0.4693.80",

	// https://blogs.opera.com/desktop/changelog-for-97/
	"111.0.5532.2,97.0.4711.0",
	"111.0.5532.2,97.0.4704.0",
	"111.0.5532.2,97.0.4697.0",
	"111.0.5562.0,97.0.4718.0",
	"111.0.5563.19,97.0.4719.4",
	"111.0.5563.19,97.0.4719.11",
	"111.0.5563.41,97.0.4719.17",
	"111.0.5563.65,97.0.4719.26",
	"111.0.5563.65,97.0.4719.28",
	"111.0.5563.111,97.0.4719.43",
	"111.0.5563.147,97.0.4719.63",
	"111.0.5563.147,97.0.4719.83",

	// https://blogs.opera.com/desktop/changelog-for-98/
	"112.0.5596.2,98.0.4756.0",
	"112.0.5596.2,98.0.4746.0",
	"112.0.5615.20,98.0.4759.1",
	"112.0.5615.50,98.0.4759.3",
	"112.0.5615.87,98.0.4759.6",
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
	"12",
	"13",
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
	"Macintosh; Intel Mac OS X 11_3",
	"Macintosh; Intel Mac OS X 11_3_1",
	"Macintosh; Intel Mac OS X 11_4",
	"Macintosh; Intel Mac OS X 11_5",
	"Macintosh; Intel Mac OS X 11_5_1",
	"Macintosh; Intel Mac OS X 11_5_2",
	"Macintosh; Intel Mac OS X 11_6",
	"Macintosh; Intel Mac OS X 11_6_1",
	"Macintosh; Intel Mac OS X 11_6_2",
	"Macintosh; Intel Mac OS X 11_6_3",
	"Macintosh; Intel Mac OS X 11_6_4",
	"Macintosh; Intel Mac OS X 11_6_5",
	"Macintosh; Intel Mac OS X 11_6_6",
	"Macintosh; Intel Mac OS X 11_6_7",
	"Macintosh; Intel Mac OS X 11_6_8",
	"Macintosh; Intel Mac OS X 11_7",
	"Macintosh; Intel Mac OS X 11_7_1",
	"Macintosh; Intel Mac OS X 11_7_2",
	"Macintosh; Intel Mac OS X 11_7_3",
	"Macintosh; Intel Mac OS X 11_7_4",
	"Macintosh; Intel Mac OS X 11_7_5",
	"Macintosh; Intel Mac OS X 11_7_6",

	// MacOS - Monterey
	"Macintosh; Intel Mac OS X 12_0",
	"Macintosh; Intel Mac OS X 12_0_1",
	"Macintosh; Intel Mac OS X 12_1",
	"Macintosh; Intel Mac OS X 12_2",
	"Macintosh; Intel Mac OS X 12_2_1",
	"Macintosh; Intel Mac OS X 12_3",
	"Macintosh; Intel Mac OS X 12_3_1",
	"Macintosh; Intel Mac OS X 12_4",
	"Macintosh; Intel Mac OS X 12_5",
	"Macintosh; Intel Mac OS X 12_5_1",
	"Macintosh; Intel Mac OS X 12_6",
	"Macintosh; Intel Mac OS X 12_6_1",
	"Macintosh; Intel Mac OS X 12_6_2",
	"Macintosh; Intel Mac OS X 12_6_3",
	"Macintosh; Intel Mac OS X 12_6_4",
	"Macintosh; Intel Mac OS X 12_6_5",

	// MacOS - Ventura
	"Macintosh; Intel Mac OS X 13_0",
	"Macintosh; Intel Mac OS X 13_0_1",
	"Macintosh; Intel Mac OS X 13_1",
	"Macintosh; Intel Mac OS X 13_2",
	"Macintosh; Intel Mac OS X 13_2_1",
	"Macintosh; Intel Mac OS X 13_3",
	"Macintosh; Intel Mac OS X 13_3_1",

	// Windows
	"Windows NT 10.0; Win64; x64",
	"Windows NT 5.1",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",

	// Linux
	"X11; Linux x86_64",
}

// Generates Firefox Browser User-Agent (Desktop)
//
//	-> "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:87.0) Gecko/20100101 Firefox/87.0"
func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

// Generates Chrome Browser User-Agent (Desktop)
//
//	-> "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"
func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

// Generates Microsoft Edge User-Agent (Desktop)
//
//	-> "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36 Edg/90.0.818.39"
func genEdgeUA() string {
	version := edgeVersions[rand.Intn(len(edgeVersions))]
	chromeVersion := strings.Split(version, ",")[0]
	edgeVersion := strings.Split(version, ",")[1]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s", os, chromeVersion, edgeVersion)
}

// Generates Opera Browser User-Agent (Desktop)
//
//	-> "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_3_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 OPR/98.0.4759.3"
func genOperaUA() string {
	version := operaVersions[rand.Intn(len(operaVersions))]
	chromeVersion := strings.Split(version, ",")[0]
	operaVersion := strings.Split(version, ",")[1]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 OPR/%s", os, chromeVersion, operaVersion)
}

// Generates UCWEB/Nokia203 Browser User-Agent (Mobile)
//
//	-> "UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCMini/10.9.8.1006 (SpeedMode; Proxy; Android 4.4.4; SM-J110H ) U2/1.0.0 Mobile"
func genMobileUcwebUA() string {
	device := ucwebDevices[rand.Intn(len(ucwebDevices))]
	version := ucwebVersions[rand.Intn(len(ucwebVersions))]
	android := androidVersions[rand.Intn(len(androidVersions))]
	return fmt.Sprintf("UCWEB/2.0 (Java; U; MIDP-2.0; Nokia203/20.37) U2/1.0.0 UCMini/%s (SpeedMode; Proxy; Android %s; %s ) U2/1.0.0 Mobile", version, android, device)
}

// Generates Nexus 10 Browser User-Agent (Mobile)
//
//	-> "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 10 Build/LMY48T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.91 Safari/537.36"
func genMobileNexus10UA() string {
	build := nexus10Builds[rand.Intn(len(nexus10Builds))]
	android := androidVersions[rand.Intn(len(androidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	safari := nexus10Safari[rand.Intn(len(nexus10Safari))]
	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Nexus 10 Build/%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/%s", android, build, chrome, safari)
}
