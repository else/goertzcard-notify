package notify

import (
	"fmt"
	"net/http"
	"net/url"
)

func (n PushoverNotifier) Notify(title, msg string) error {
	values := url.Values{
		"user":    []string{n.User},
		"token":   []string{n.Token},
		"title":   []string{title},
		"message": []string{msg},
	}

	resp, err := http.PostForm("https://api.pushover.net/1/messages.json", values)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received invalid status code: %s", resp.Status)
	}

	return nil
}
