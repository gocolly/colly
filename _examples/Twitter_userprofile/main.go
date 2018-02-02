package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/gocolly/colly"
)


// Userprofile store twitter user data
type Userprofile struct {
	FullName string
	UserName string
	AccountVerified string
	Bio string
	Location string
	Website string
	JoinDate string
	Tweets string
	Following string
	Followers string
	Likes string
	Lists string
	Moments string
}

var (
	tempuserprofile = Userprofile{};
)
func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Cache responses to prevent multiple download of pages
		// even if the collector is restarted
		//colly.CacheDir("./cache"),
	)

	userprofiles := make([]Userprofile, 0, 1000)
	// On ProfileHeaderCard class
	c.OnHTML(".ProfileHeaderCard", func(e *colly.HTMLElement) {
			tempuserprofile.FullName = e.ChildText(".ProfileHeaderCard-nameLink")
      tempuserprofile.UserName = e.ChildText(".ProfileHeaderCard-screenname")
			tempuserprofile.AccountVerified = e.ChildText("div#page-container > div:nth-of-type(2) > div > div > div > div > div > div > div > h1 > span > a > span > span")
			tempuserprofile.Bio = e.ChildText(".ProfileHeaderCard-bio")
			tempuserprofile.Location = e.ChildText(".ProfileHeaderCard-locationText")
			tempuserprofile.Website = e.ChildText("div.ProfileHeaderCard-url")
			tempuserprofile.JoinDate = e.ChildAttr(".ProfileHeaderCard-joinDateText","title")
	})

		c.OnHTML("div.ProfileNav", func(e *colly.HTMLElement) {
			tempuserprofile.Tweets = e.ChildAttr("ul.ProfileNav-list > li.ProfileNav-item--tweets > a > span:nth-of-type(3)","data-count")
			if tempuserprofile.Tweets == "" {
				tempuserprofile.Tweets = "No Tweets found"
			}
			tempuserprofile.Following = e.ChildAttr("ul.ProfileNav-list > li.ProfileNav-item--following > a > span:nth-of-type(3)","data-count")
			if tempuserprofile.Following == "" {
				tempuserprofile.Following = "No Following found"
			}
			tempuserprofile.Followers = e.ChildAttr("ul.ProfileNav-list > li.ProfileNav-item--followers > a > span:nth-of-type(3)","data-count")
			if tempuserprofile.Followers == "" {
				tempuserprofile.Followers = "No Followers found"
			}
			tempuserprofile.Likes = e.ChildAttr("ul.ProfileNav-list > li.ProfileNav-item--favorites > a > span:nth-of-type(3)","data-count")
			if tempuserprofile.Likes == "" {
				tempuserprofile.Likes = "No Likes found"
			}
			tempuserprofile.Lists = e.ChildText("ul.ProfileNav-list > li.ProfileNav-item--lists > a > span:nth-of-type(3)")
			if tempuserprofile.Lists == "" {
				tempuserprofile.Lists = "No Lists found"
			}
			tempuserprofile.Moments = e.ChildText("ul.ProfileNav-list > li.ProfileNav-item--moments > a > span:nth-of-type(3)")
			if tempuserprofile.Moments == "" {
				tempuserprofile.Moments = "No Moments found"
			}
			userprofiles = append(userprofiles, tempuserprofile)
			fmt.Println(tempuserprofile)
			//more url here
			//c.Visit("https://twitter.com/" + "bharatsewani199")

})

	// Start scraping on https://twitter.com/BetaList
	c.Visit("https://twitter.com/" + "bharatsewani199")
	//c.Visit("https://twitter.com/BetaList")

	// Convert results to JSON data if the scraping job has finished
	jsonData, err := json.MarshalIndent(userprofiles, "", "  ")
	if err != nil {
		panic(err)
	}

	// Dump json to the standard output (can be redirected to a file)
	fmt.Println(string(jsonData))
	err = ioutil.WriteFile("output.json", jsonData, 0644)
}
