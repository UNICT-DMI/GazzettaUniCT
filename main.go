package main

import (
	"GazzettaUniCT/config"
	"GazzettaUniCT/history"
	"GazzettaUniCT/telegram"
	"fmt"
	"log"
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

func createHistoryFileIfNotExist() {
	_, err := os.Stat("data/history.json")

	if os.IsNotExist(err) {
		f, errFile := os.Create("data/history.json")
		if errFile != nil {
			log.Fatal(err)
		}

		f.Close()
	}
}

func main() {
	createDataFolderIfNotExist()
	createHistoryFileIfNotExist()
	const historyPath = "./data/history.json"
	const url = "http://www.oocc.unict.it/oocc/vis_verb.asp?oocc=2"

	conf, err := config.LoadConfig()

	fmt.Println(conf.BotApiKey + " " + conf.ChannelName)

	if err != nil {
		log.Panic(err)
	}

	// Instantiate default collector
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		extensionFile := strings.ToLower(link[len(link)-4:])

		if extensionFile == ".pdf" && link[2:9] == "verbali" {
			if len(link) > 45 {
				documentType := strings.ToLower(link[10:17])

				if documentType == "verbale" {
					fileUrl := "http://www.oocc.unict.it/oocc" + link[1:]

					if history.EntryNotExist(fileUrl, historyPath) {
						err := telegram.SendDocument(conf.BotApiKey, conf.ChannelName, fileUrl)

						if err != nil {
							fmt.Println("Error sending message: " + fileUrl)
						} else {
							err := history.NewHistoryEntry(fileUrl, historyPath)
							log.Println("Message " + path.Base(fileUrl) + " sended!")

							if err != nil {
								log.Println(err)
							}
						}
					}
				}
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
