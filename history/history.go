package history

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

type historyEntry struct {
	Name string `json:"name"`
}

func getNameFromUrl(url string) string {
	name := path.Base(url)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func NewHistoryEntry(url string, pathfile string) error {
	data, err := ioutil.ReadFile(pathfile)
	name := getNameFromUrl(url)

	var history []historyEntry

	if err := json.Unmarshal(data, &history); err != nil {
		log.Println(err)
	}

	history = append(history, historyEntry{name})

	dataBytes, err := json.Marshal(history)

	err = ioutil.WriteFile(pathfile, dataBytes, 0644)

	return err
}

func EntryNotExist(url string, pathfile string) bool {
	data, err := ioutil.ReadFile(pathfile)
	name := getNameFromUrl(url)

	if err != nil {
		log.Panic(err)
	}

	var history []historyEntry

	if err := json.Unmarshal(data, &history); err != nil {
		log.Println(err)
	}

	for _, entry := range history {
		if entry.Name == name {
			return false
		}
	}

	return true
}
