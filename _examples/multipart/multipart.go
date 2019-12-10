package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

func generateFormData() map[string][]byte {
	f, _ := os.Open("gocolly.jpg")
	defer f.Close()

	imgData, _ := ioutil.ReadAll(f)

	return map[string][]byte{
		"firstname": []byte("one"),
		"lastname":  []byte("two"),
		"email":     []byte("onetwo@example.com"),
		"file":      imgData,
	}
}

func setupServer() {
	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received request")
		err := r.ParseMultipartForm(10000000)
		if err != nil {
			fmt.Println("server: Error")
			w.WriteHeader(500)
			w.Write([]byte("<html><body>Internal Server Error</body></html>"))
			return
		}
		w.WriteHeader(200)
		fmt.Println("server: OK")
		w.Write([]byte("<html><body>Success</body></html>"))
	}

	go http.ListenAndServe(":8080", handler)
}

func main() {
	// Start a single route http server to post an image to.
	setupServer()

	c := colly.NewCollector(colly.AllowURLRevisit(), colly.MaxDepth(5))

	// On every a element which has href attribute call callback
	c.OnHTML("html", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		time.Sleep(1 * time.Second)
		e.Request.PostMultipart("http://localhost:8080/", generateFormData())
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Posting gocolly.jpg to", r.URL.String())
	})

	// Start scraping
	c.PostMultipart("http://localhost:8080/", generateFormData())
	c.Wait()
}
