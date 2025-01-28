package notification

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

const notificationUrl string = "https://ntfy.sh/"

type Message struct {
	Title string
	Body  string
	Tags  string
}

func Send(msg Message) error {
	ntfyKey := os.Getenv("BIN_JUICE_NTFY_KEY")
	if ntfyKey == "" {
		return fmt.Errorf("error: no notification key found")
	}

	req, _ := http.NewRequest("POST", notificationUrl+ntfyKey, strings.NewReader(msg.Body))
	req.Header.Set("Title", msg.Title)
	req.Header.Set("Tags", msg.Tags)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error: failed to send notification: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error: status code %d response from notification", res.StatusCode)
	}
	return nil
}
