package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/gocolly/colly/v2/queue"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	IstockSearchAPI        = "https://www.istockphoto.com/search/2/image"
	ColorSimilarityAssetid = "colorsimilarityassetid"

	MaxPages       = 20
	MinPages       = 1
	MaxPower       = 64
	MinPower       = 1
	DefaultBackend = "istock_dataset"

	Content = "content"
	Color   = "color"
)

type Downloader struct {
	// Phrase is the image tag keyword to be retrieved
	Phrase string
	// Pages is the size of the data that needs to be collected
	// For demonstration purposes, don't let Pages exceed MinPages and MaxPages
	// During initialization, invalid Pages values will be automatically corrected
	Pages int
	// MediaType defaults to PHOTO, options can be viewed in typing
	Mediatype string
	// NumberOfPeople defaults to NoPeople, options can be viewed in typing
	NumberOfPeople string
	// Orientations defaults to SQUARE, options can be viewed in typing
	Orientations string
	// Backend is the root directory of the image cache
	// the default value is DefaultBackend
	Backend string
	// Flag is the name of the parent directory where images are stored,
	// and its default value is the keyword you specify, namely Phrase
	Flag     string
	Power    int
	Similar  string
	ProxyURL string

	dirLocal string
	holdAPI  string
	query    string

	collector *colly.Collector
	worker    *queue.Queue
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// NewDownloader Initialize the downloader object
func NewDownloader(phrase string) *Downloader {
	phrase = strings.Trim(phrase, " ")
	if phrase == "" {
		log.Fatalln("Invalid phrase")
	}

	d := &Downloader{Phrase: phrase}
	d.Init()
	return d
}

func (d *Downloader) Init() {
	d.Mediatype = queryDefault[MediaType]
	d.NumberOfPeople = queryDefault[NumberOfPeople]
	d.Orientations = queryDefault[Orientations]
	d.Flag = d.Phrase
	d.Pages = MinPages
	d.Backend = DefaultBackend
	d.Power = runtime.NumCPU()
	d.holdAPI = IstockSearchAPI
	d.Similar = Content

	d.collector = colly.NewCollector()
	d.worker, _ = queue.New(1, nil)
}

// MoreLikeThis Similarity search
func (d *Downloader) MoreLikeThis(istockID int) *Downloader {
	var similarMatch = map[string]string{
		Content: fmt.Sprintf("https://www.istockphoto.com/search/more-like-this/%d", istockID),
		Color:   fmt.Sprintf("https://www.istockphoto.com/search/2/image?%s=%d", ColorSimilarityAssetid, istockID),
	}
	d.holdAPI = similarMatch[d.Similar]

	return d
}

// Mining Start the collector
func (d *Downloader) Mining() {
	d.preload()
	d.overload()

	if err := d.worker.Run(d.collector); err != nil {
		log.Fatalln("Failed to setup worker, ", err)
	}
	log.Println("Task complete.")
}

func (d *Downloader) preload() {
	d.checkParams()
	d.checkWorkspace()
	d.checkQuery()
	d.initWorker()

	log.Printf("Container preload - phrase=`%s`", d.Phrase)
	log.Printf("Setup [istock] - power=%d pages=%d", d.Power, d.Pages)
}

func (d *Downloader) checkParams() {
	if d.Pages > MaxPages || d.Pages < 1 {
		log.Printf("Automatically calibrate to default values. - pages∈[%d, %d]\n", MinPages, MaxPages)
		d.Pages = MinPages
	}

	d.Mediatype = RefactorInvalidQueryType(MediaType, d.Mediatype)
	d.Orientations = RefactorInvalidQueryType(Orientations, d.Orientations)
	d.NumberOfPeople = RefactorInvalidQueryType(NumberOfPeople, d.NumberOfPeople)
}

func (d *Downloader) checkWorkspace() {
	var badCode = []string{"\\", "/", ":", "*", "?", "\"", "<", ">", "|", " ", "."}

	for _, c := range badCode {
		strings.ReplaceAll(c, d.Flag, d.Flag)
	}

	if d.Backend == DefaultBackend {
		d.dirLocal = path.Join(d.Backend, d.Flag)
	} else {
		d.dirLocal = path.Join(d.Backend, DefaultBackend, d.Flag)
	}

	err := os.MkdirAll(d.dirLocal, os.ModePerm)
	if err != nil {
		log.Fatalln("WorkspaceCheckerException: ", err)
	}
}

func (d *Downloader) checkQuery() {
	var params string
	parser, _ := url.Parse(d.holdAPI)
	if parser.Path == "/search/2/image" && strings.HasPrefix(parser.RawQuery, ColorSimilarityAssetid) {
		params = fmt.Sprintf("%s&phrase=%s", d.holdAPI, d.Phrase)
	} else {
		params = fmt.Sprintf("%s?phrase=%s", d.holdAPI, d.Phrase)
	}

	if d.Mediatype != UNDEFINED {
		params += fmt.Sprintf("&mediatype=%s", d.Mediatype)
	}
	if d.NumberOfPeople != UNDEFINED {
		params += fmt.Sprintf("&numberofpeople=%s", d.NumberOfPeople)
	}
	if d.Orientations != UNDEFINED {
		params += fmt.Sprintf("&orientations=%s", d.Orientations)
	}

	d.query = params
}

func (d *Downloader) initWorker() {
	// [1] Init concurrent-tasks
	for i := 1; i < d.Pages+1; i++ {
		URL := fmt.Sprintf("%s&page=%d", d.query, i)
		err := d.worker.AddURL(URL)
		if err != nil {
			log.Fatalln("DownloaderPreloadException: ", err)
		} else {
			log.Println("SetEntrance: ", URL)
		}
	}

	// [2] Reset threads of the worker
	if d.Power > MaxPower || d.Power < MinPower {
		log.Printf("Automatically calibrate to default values. - power∈[%d, %d]\n", MinPower, MaxPower)
		d.Power = MinPower
	}
	if d.Power >= d.Pages {
		d.Power = d.Pages
	}
	d.worker.Threads = d.Power

	// [3] Refactor Colly Headers
	extensions.Referer(d.collector)
	d.collector.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
		"Chrome/103.0.5060.134 Safari/537.36 Edg/103.0.1264.77"

	// CN：这是一个被墙掉的网站，必须使用代理访问，参考代码如下：
	// d := NewDownloader("phrase")
	// d.ProxyURL = "http://127.0.0.1:10809"
	// d.Mining()
	if d.ProxyURL != "" {
		if err := d.collector.SetProxy(d.ProxyURL); err != nil {
			log.Printf("Failed to set collector's proxy - err=%s", err)
		}
	}

}

func (d *Downloader) overload() {
	d.collector.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 0 {
			log.Println("HTTPConnectionError:", err)
		} else {
			log.Println(err)
		}
	})

	d.collector.OnHTML("img.MosaicAsset-module__thumb___klD9E", func(e *colly.HTMLElement) {
		//TODO Extract istock ID, remove duplicate tasks
		imageURL := e.Attr("src")
		if err := d.worker.AddURL(imageURL); err != nil {
			log.Printf("Failed to download image - URL=%s", imageURL)
		}

	})

	d.collector.OnScraped(func(r *colly.Response) {
		if progress, _ := d.worker.Size(); progress != 0 {
			log.Printf("Offload - taskID=%s progess=%d", r.FileName(), progress)
		}
		if filepath.Ext(r.FileName()) == ".jpg" {
			fn := path.Join(d.dirLocal, r.FileName())
			if err := r.Save(fn); err != nil {
				log.Printf("Failed to offload - URL=%s", r.Request.URL.String())
			}
		}
	})
}
