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
	"github.com/gocolly/colly/queue"
)

const historyPath = "./data/history.json"
const urlConsiglioAmministrazione = "http://www.oocc.unict.it/oocc/vis_verb.asp?oocc=1"
const urlSenato = "http://www.oocc.unict.it/oocc/vis_verb.asp?oocc=2"

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

	conf, err := config.LoadConfig()

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

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	q.AddURL(urlSenato)
	q.AddURL(urlConsiglioAmministrazione)

	// Start scraping
	q.Run(c)
}
