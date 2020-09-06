package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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

func DownloadFile(url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Get the filename
	fileName := path.Base(url)
	filePath := "./data/" + fileName

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
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

			err := DownloadFile(fileUrl)

			if err != nil {
				fmt.Println("Error download: " + fileUrl)
			}

		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping
	c.Visit(url)
}
