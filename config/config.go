package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	BotApiKey   string `json:"bot_api_key"`
	ChannelName string `json:"channel_name"`
}

const configPath = "data/config/config.json"

func LoadConfig() (Config, error) {
	data, err := ioutil.ReadFile(configPath)
	var config Config

	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	}

	return config, err
}
