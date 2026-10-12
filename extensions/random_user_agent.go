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
	genMobilePixel7UA,
	genMobilePixel6UA,
	genMobilePixel5UA,
	genMobilePixel4UA,
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
	113.0,
	114.0,
	115.0,
	116.0,
	117.0,
	118.0,
	119.0,
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

	// https://chromereleases.googleblog.com/2023/05/stable-channel-update-for-desktop.html
	"113.0.5672.63",
	"113.0.5672.64",

	// https://chromereleases.googleblog.com/2023/05/stable-channel-update-for-desktop_8.html
	"113.0.5672.92",
	"113.0.5672.93",

	// https://chromereleases.googleblog.com/2023/05/stable-channel-update-for-desktop_16.html
	"113.0.5672.126",
	"113.0.5672.127",

	// https://chromereleases.googleblog.com/2023/05/stable-channel-update-for-desktop_30.html
	"114.0.5735.90/91",
	"114.0.5735.91",

	// https://chromereleases.googleblog.com/2023/06/stable-channel-update-for-desktop.html
	"114.0.5735.106",
	"114.0.5735.110",

	// https://chromereleases.googleblog.com/2023/06/stable-channel-update-for-desktop_13.html
	"114.0.5735.133",
	"114.0.5735.134",

	// https://chromereleases.googleblog.com/2023/06/stable-channel-update-for-desktop_26.html
	"114.0.5735.198",
	"114.0.5735.199",

	// https://chromereleases.googleblog.com/2023/07/stable-channel-update-for-desktop.html
	"115.0.5790.98",
	"115.0.5790.99",

	// https://chromereleases.googleblog.com/2023/07/stable-channel-update-for-desktop_20.html
	"115.0.5790.102",

	// https://chromereleases.googleblog.com/2023/07/stable-channel-update-for-desktop_25.html
	"115.0.5790.110",
	"115.0.5790.114",

	// https://chromereleases.googleblog.com/2023/08/stable-channel-update-for-desktop.html
	"115.0.5790.170",
	"115.0.5790.171",

	// https://chromereleases.googleblog.com/2023/08/stable-channel-update-for-desktop_15.html
	"116.0.5845.96",
	"116.0.5845.97",

	// https://chromereleases.googleblog.com/2023/08/chrome-desktop-stable-update.html
	"116.0.5845.110",
	"116.0.5845.111",

	// https://chromereleases.googleblog.com/2023/08/stable-channel-update-for-desktop_29.html
	"116.0.5845.140",
	"116.0.5845.141",

	// https://chromereleases.googleblog.com/2023/09/stable-channel-update-for-desktop.html
	"116.0.5845.179",
	"116.0.5845.180",

	// https://chromereleases.googleblog.com/2023/09/stable-channel-update-for-desktop_11.html
	"116.0.5845.187",
	"116.0.5845.188",

	// https://chromereleases.googleblog.com/2023/09/stable-channel-update-for-desktop_12.html
	"117.0.5938.62",
	"117.0.5938.63",

	// https://chromereleases.googleblog.com/2023/09/stable-channel-update-for-desktop_15.html
	"117.0.5938.88",
	"117.0.5938.89",

	// https://chromereleases.googleblog.com/2023/09/stable-channel-update-for-desktop_21.html
	"117.0.5938.92",

	// https://chromereleases.googleblog.com/2023/09/stable-channel-update-for-desktop_27.html
	"117.0.5938.132",

	// https://chromereleases.googleblog.com/2023/10/stable-channel-update-for-desktop.html
	"117.0.5938.149",
	"117.0.5938.150",

	// https://chromereleases.googleblog.com/2023/10/stable-channel-update-for-desktop_10.html
	"118.0.5993.70",
	"118.0.5993.71",

	// https://chromereleases.googleblog.com/2023/10/stable-channel-update-for-desktop_17.html
	"118.0.5993.88",
	"118.0.5993.89",

	// https://chromereleases.googleblog.com/2023/10/stable-channel-update-for-desktop_19.html
	"118.0.5993.96",

	// https://chromereleases.googleblog.com/2023/10/stable-channel-update-for-desktop_24.html
	"118.0.5993.117",
	"118.0.5993.118",

	// https://chromereleases.googleblog.com/2023/10/stable-channel-update-for-desktop_31.html
	"119.0.6045.105",
	"119.0.6045.106",

	// https://chromereleases.googleblog.com/2023/11/stable-channel-update-for-desktop.html
	"119.0.6045.123",
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
	"113.0.0.0,113.0.1774.3",
	"114.0.0.0,114.0.1823.37",
	"115.0.0.0,115.0.1901.183",
	"116.0.0.0,116.0.1938.54",
	"117.0.0.0,117.0.2045.31",
	"118.0.0.0,118.0.2088.46",
	"119.0.0.0,119.0.2151.44",
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
	"112.0.5615.165,98.0.4759.15",
	"112.0.5615.165,98.0.4759.21",
	"112.0.5615.165,98.0.4759.39",

	// https://blogs.opera.com/desktop/changelog-for-99/
	"113.0.5672.64,99.0.4788.5",
	"113.0.5672.93,99.0.4788.9",
	"113.0.5672.127,99.0.4788.13",
	"113.0.5672.127,99.0.4788.31",
	"113.0.5672.127,99.0.4788.47",
	"113.0.5672.127,99.0.4788.65",
	"113.0.5672.127,99.0.4788.77",
	"113.0.5672.127,99.0.4788.88",

	// https://blogs.opera.com/desktop/changelog-for-100/
	"114.0.5720.4,100.0.4809.0",
	"114.0.5735.9,100.0.4815.0",
	"114.0.5735.9,100.0.4815.2",
	"114.0.5735.110,100.0.4815.13",
	"114.0.5735.199,100.0.4815.30",
	"114.0.5735.199,100.0.4815.47",
	"114.0.5735.199,100.0.4815.54",
	"114.0.5735.199,100.0.4815.76",

	// https://blogs.opera.com/desktop/changelog-for-101/
	"115.0.5790.3,101.0.4843.0",
	"115.0.5790.3,101.0.4843.5",
	"115.0.5790.40,101.0.4843.10",
	"115.0.5790.40,101.0.4843.13",
	"115.0.5790.75,101.0.4843.19",
	"115.0.5790.102,101.0.4843.25",
	"115.0.5790.102,101.0.4843.33",
	"115.0.5790.171,101.0.4843.43",
	"115.0.5790.171,101.0.4843.58",

	// https://blogs.opera.com/desktop/changelog-for-102/
	"116.0.5829.0,102.0.4871.0",
	"116.0.5845.42,102.0.4879.0",
	"116.0.5845.42,102.0.4880.6",
	"116.0.5845.62,102.0.4880.10",
	"116.0.5845.97,102.0.4880.16",
	"116.0.5845.97,102.0.4880.28",
	"116.0.5845.141,102.0.4880.33",
	"116.0.5845.141,102.0.4880.38",
	"116.0.5845.141,102.0.4880.40",
	"116.0.5845.141,102.0.4880.46",
	"116.0.5845.141,102.0.4880.51",
	"116.0.5845.141,102.0.4880.56",
	"116.0.5845.141,102.0.4880.70",
	"116.0.5845.141,102.0.4880.78",

	// https://blogs.opera.com/desktop/changelog-for-103/
	"117.0.5897.3,103.0.4920.0",
	"117.0.5938.0,103.0.4928.0",
	"117.0.5938.132,103.0.4928.16",
	"117.0.5938.132,103.0.4928.26",
	"117.0.5938.132,103.0.4928.34",

	// https://blogs.opera.com/desktop/changelog-for-104/
	"118.0.5993.11,104.0.4944.3",
	"118.0.5993.11,104.0.4944.10",
	"118.0.5993.71,104.0.4944.18",
	"118.0.5993.71,104.0.4944.23",
	"118.0.5993.71,104.0.4944.28",
	"118.0.5993.96,104.0.4944.33",
	"118.0.5993.118,104.0.4944.36",
	"118.0.5993.118,104.0.4944.54",

	// https://blogs.opera.com/desktop/changelog-for-105/
	"119.0.6034.6,105.0.4963.0",
	"119.0.6045.9,105.0.4970.6",
	"119.0.6045.105,105.0.4970.10",
}

var pixel7AndroidVersions = []string{
	// Data source:
	// - https://developer.android.com/about/versions
	// - https://source.android.com/docs/setup/about/build-numbers#source-code-tags-and-builds
	"13",
}

var pixel6AndroidVersions = []string{
	// Data source:
	// - https://developer.android.com/about/versions
	// - https://source.android.com/docs/setup/about/build-numbers#source-code-tags-and-builds
	"12",
	"13",
}

var pixel5AndroidVersions = []string{
	// Data source:
	// - https://developer.android.com/about/versions
	// - https://source.android.com/docs/setup/about/build-numbers#source-code-tags-and-builds
	"11",
	"12",
	"13",
}

var pixel4AndroidVersions = []string{
	// Data source:
	// - https://developer.android.com/about/versions
	// - https://source.android.com/docs/setup/about/build-numbers#source-code-tags-and-builds
	"10",
	"11",
	"12",
	"13",
}

var nexus10AndroidVersions = []string{
	// Data source:
	// - https://developer.android.com/about/versions
	// - https://source.android.com/docs/setup/about/build-numbers#source-code-tags-and-builds
	"4.4.2",
	"4.4.4",
	"5.0",
	"5.0.1",
	"5.0.2",
	"5.1",
	"5.1.1",
}

var nexus10Builds = []string{
	// Data source: https://source.android.com/docs/setup/about/build-numbers#source-code-tags-and-builds

	"LMY49M", // android-5.1.1_r38 (Lollipop)
	"LMY49J", // android-5.1.1_r37 (Lollipop)
	"LMY49I", // android-5.1.1_r36 (Lollipop)
	"LMY49H", // android-5.1.1_r35 (Lollipop)
	"LMY49G", // android-5.1.1_r34 (Lollipop)
	"LMY49F", // android-5.1.1_r33 (Lollipop)
	"LMY48Z", // android-5.1.1_r30 (Lollipop)
	"LMY48X", // android-5.1.1_r25 (Lollipop)
	"LMY48T", // android-5.1.1_r19 (Lollipop)
	"LMY48M", // android-5.1.1_r14 (Lollipop)
	"LMY48I", // android-5.1.1_r9 (Lollipop)
	"LMY47V", // android-5.1.1_r1 (Lollipop)
	"LMY47D", // android-5.1.0_r1 (Lollipop)
	"LRX22G", // android-5.0.2_r1 (Lollipop)
	"LRX22C", // android-5.0.1_r1 (Lollipop)
	"LRX21P", // android-5.0.0_r4.0.1 (Lollipop)
	"KTU84P", // android-4.4.4_r1 (KitKat)
	"KTU84L", // android-4.4.3_r1 (KitKat)
	"KOT49H", // android-4.4.2_r1 (KitKat)
	"KOT49E", // android-4.4.1_r1 (KitKat)
	"KRT16S", // android-4.4_r1.2 (KitKat)
	"JWR66Y", // android-4.3_r1.1 (Jelly Bean)
	"JWR66V", // android-4.3_r1 (Jelly Bean)
	"JWR66N", // android-4.3_r0.9.1 (Jelly Bean)
	"JDQ39 ", // android-4.2.2_r1 (Jelly Bean)
	"JOP40F", // android-4.2.1_r1.1 (Jelly Bean)
	"JOP40D", // android-4.2.1_r1 (Jelly Bean)
	"JOP40C", // android-4.2_r1 (Jelly Bean)
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
// -> "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:87.0) Gecko/20100101 Firefox/87.0"
func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

// Generates Chrome Browser User-Agent (Desktop)
//
// -> "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36"
func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

// Generates Microsoft Edge User-Agent (Desktop)
//
// -> "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.72 Safari/537.36 Edg/90.0.818.39"
func genEdgeUA() string {
	version := edgeVersions[rand.Intn(len(edgeVersions))]
	chromeVersion := strings.Split(version, ",")[0]
	edgeVersion := strings.Split(version, ",")[1]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 Edg/%s", os, chromeVersion, edgeVersion)
}

// Generates Opera Browser User-Agent (Desktop)
//
// -> "Mozilla/5.0 (Macintosh; Intel Mac OS X 13_3_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 OPR/98.0.4759.3"
func genOperaUA() string {
	version := operaVersions[rand.Intn(len(operaVersions))]
	chromeVersion := strings.Split(version, ",")[0]
	operaVersion := strings.Split(version, ",")[1]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36 OPR/%s", os, chromeVersion, operaVersion)
}

// Generates Pixel 7 Browser User-Agent (Mobile)
//
// -> Mozilla/5.0 (Linux; Android 13; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36
func genMobilePixel7UA() string {
	android := pixel7AndroidVersions[rand.Intn(len(pixel7AndroidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Pixel 7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", android, chrome)
}

// Generates Pixel 6 Browser User-Agent (Mobile)
//
// -> "Mozilla/5.0 (Linux; Android 13; Pixel 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36"
func genMobilePixel6UA() string {
	android := pixel6AndroidVersions[rand.Intn(len(pixel6AndroidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Pixel 6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", android, chrome)
}

// Generates Pixel 5 Browser User-Agent (Mobile)
//
// -> "Mozilla/5.0 (Linux; Android 13; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36"
func genMobilePixel5UA() string {
	android := pixel5AndroidVersions[rand.Intn(len(pixel5AndroidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", android, chrome)
}

// Generates Pixel 4 Browser User-Agent (Mobile)
//
// -> "Mozilla/5.0 (Linux; Android 13; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36"
func genMobilePixel4UA() string {
	android := pixel4AndroidVersions[rand.Intn(len(pixel4AndroidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", android, chrome)
}

// Generates Nexus 10 Browser User-Agent (Mobile)
//
// -> "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 10 Build/LMY48T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.91 Safari/537.36"
func genMobileNexus10UA() string {
	build := nexus10Builds[rand.Intn(len(nexus10Builds))]
	android := nexus10AndroidVersions[rand.Intn(len(nexus10AndroidVersions))]
	chrome := chromeVersions[rand.Intn(len(chromeVersions))]
	return fmt.Sprintf("Mozilla/5.0 (Linux; Android %s; Nexus 10 Build/%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", android, build, chrome)
}
