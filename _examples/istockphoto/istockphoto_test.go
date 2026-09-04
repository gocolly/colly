package main

import (
	"fmt"
	"github.com/gocolly/colly/v2/_examples/istockphoto/downloader"
	"testing"
)

const (
	// Developers in mainland China need to use a proxy for testing
	// Because `istockphoto.com` is a website blocked by GFW
	// Sample: ProxyURL = "http://127.0.0.1:10809"
	// IF ProxyURL == "" then this configuration will not take place
	ProxyURL = ""
)

func TestStandardDownloader(t *testing.T) {
	downloader.NewDownloader("cyberpunk").Mining()
}

func TestDownloaderWithProxyURL(t *testing.T) {
	d := downloader.NewDownloader("cyberpunk")
	d.ProxyURL = ProxyURL
	d.Mining()
}

func TestDownloaderNotQuery(t *testing.T) {
	d := downloader.NewDownloader("dog")
	d.ProxyURL = ProxyURL
	d.CloseFilter()
	d.Mining()
}

func TestDownloaderWithPages(t *testing.T) {
	d := downloader.NewDownloader("cat")
	d.Pages = 4
	d.ProxyURL = ProxyURL
	d.Mining()
}

func TestDownloaderWithFlag(t *testing.T) {
	// Images with different label will be centrally stored in the same directory
	flag := "lion"
	phrases := []string{"lion closed eyes", "lion open mouth"}

	for _, phrase := range phrases {
		d := downloader.NewDownloader(phrase)
		d.Orientations = downloader.Orientations.Undefined
		d.Flag = flag
		d.ProxyURL = ProxyURL
		d.Mining()
	}
}

func TestMoreLikeThisContent(t *testing.T) {
	d := downloader.NewDownloader("gun")
	d.Orientations = downloader.Orientations.Undefined
	d.Flag = "gun-similar-content"
	d.ProxyURL = ProxyURL
	d.MoreLikeThis(529989264).Mining()
}

func TestMoreLikeThisColor(t *testing.T) {
	tag := "cyberpunk"
	d := downloader.NewDownloader(tag)
	d.Orientations = downloader.Orientations.Undefined
	d.Flag = fmt.Sprintf("%s-similar-color", tag)
	d.Similar = downloader.Color
	d.ProxyURL = ProxyURL
	d.MoreLikeThis(1266931346).Mining()
}
