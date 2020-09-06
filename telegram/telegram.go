package telegram

import (
	"fmt"
	"net/http"
)

func SendDocument(botApiKey string, channelName string, url string) error {
	urlRequest := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%s&photo=https://cutt.ly/gazzetta_unict_test_photo&reply_&reply_markup={\"inline_keyboard\":[[{\"text\": \"Vai al verbale\", \"url\": \"%s\"}]]}",
		botApiKey, channelName, url)

	_, err := http.Get(urlRequest)

	return err
}
