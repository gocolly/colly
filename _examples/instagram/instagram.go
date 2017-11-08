package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asciimoo/colly"
)

// found in https://www.instagram.com/static/bundles/en_US_Commons.js/68e7390c5938.js
// included from profile page
const instagramQueryId string = "17888483320059182"

// "id": user id, "after": end cursor
const nextPageURLTemplate string = `https://www.instagram.com/graphql/query/?query_id=17888483320059182&variables={"id":"%s","first":12,"after":"%s"}`

type pageInfo struct {
	EndCursor string `json:"end_cursor"`
	NextPage  bool   `json:"has_next_page"`
}

func main() {
	if len(os.Args) != 2 {
		log.Println("Missing account name argument")
		os.Exit(1)
	}

	var actualUserId string
	instagramAccount := os.Args[1]
	outputDir := fmt.Sprintf("./instagram_%s/", instagramAccount)

	c := colly.NewCollector()
	c.CacheDir = "./_instagram_cache/"
	c.UserAgent = "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"

	c.OnHTML("body > script:first-of-type", func(e *colly.HTMLElement) {
		jsonData := e.Text[strings.Index(e.Text, "{") : len(e.Text)-1]
		data := struct {
			EntryData struct {
				ProfilePage []struct {
					User struct {
						Id    string `json:"id"`
						Media struct {
							Nodes []struct {
								ImageURL     string `json:"display_src"`
								ThumbnailURL string `json:"thumbnail_src"`
								IsVideo      bool   `json:"is_video"`
								Date         int    `json:"date"`
								Dimensions   struct {
									Width  int `json:"width"`
									Height int `json:"height"`
								}
							}
							PageInfo pageInfo `json:"page_info"`
						} `json:"media"`
					} `json:"user"`
				} `json:"ProfilePage"`
			} `json:"entry_data"`
		}{}
		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("saving output to ", outputDir)
		os.MkdirAll(outputDir, os.ModePerm)
		page := data.EntryData.ProfilePage[0]
		actualUserId = page.User.Id
		for _, obj := range page.User.Media.Nodes {
			// skip videos
			if obj.IsVideo {
				continue
			}
			c.Visit(obj.ImageURL)
		}
		if page.User.Media.PageInfo.NextPage {
			log.Println("Next page found")
			c.Visit(fmt.Sprintf(nextPageURLTemplate, actualUserId, page.User.Media.PageInfo.EndCursor))
		}
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Index(r.Headers.Get("Content-Type"), "image") > -1 {
			r.Save(outputDir + r.FileName())
			return
		}

		if strings.Index(r.Headers.Get("Content-Type"), "json") == -1 {
			return
		}

		data := struct {
			Data struct {
				User struct {
					Container struct {
						PageInfo pageInfo `json:"page_info"`
						Edges    []struct {
							Node struct {
								ImageURL     string `json:"display_url"`
								ThumbnailURL string `json:"thumbnail_src"`
								IsVideo      bool   `json:"is_video"`
								Date         int    `json:"taken_at_timestamp"`
								Dimensions   struct {
									Width  int `json:"width"`
									Height int `json:"height"`
								}
							}
						} `json:"edges"`
					} `json:"edge_owner_to_timeline_media"`
				}
			} `json:"data"`
		}{}
		err := json.Unmarshal(r.Body, &data)
		if err != nil {
			log.Fatal(err)
		}

		for _, obj := range data.Data.User.Container.Edges {
			// skip videos
			if obj.Node.IsVideo {
				continue
			}
			c.Visit(obj.Node.ImageURL)
		}
		if data.Data.User.Container.PageInfo.NextPage {
			log.Println("Next page found")
			c.Visit(fmt.Sprintf(nextPageURLTemplate, actualUserId, data.Data.User.Container.PageInfo.EndCursor))
		}
	})

	c.Visit("https://instagram.com/" + instagramAccount)
}
