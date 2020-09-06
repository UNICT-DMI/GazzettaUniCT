package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func createDataFolderIfNotExist() {
	_, err := os.Stat("data")

	if os.IsNotExist(err) {
		errDir := os.MkdirAll("data", 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	createDataFolderIfNotExist()

	// Instantiate default collector
	c := colly.NewCollector()
	url := "http://www.oocc.unict.it/oocc/vis_verb.asp?oocc=2"

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		extensionFile := strings.ToLower(link[len(link)-4:])

		if extensionFile == ".pdf" && link[2:9] == "verbali" {
			fileUrl := "http://www.oocc.unict.it/oocc" + link[1:]
			fmt.Println("fileUrl: " + fileUrl)
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	c.Visit(url)
}
