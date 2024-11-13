package notify

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"

	"github.com/nice-pink/goutil/pkg/log"
)

type SlackMessage struct {
	Url   string
	Text  string
	Info  string
	Color string
	// Status string
}

func Send(msg SlackMessage) error {
	attachment := GetAttachment(msg)
	body := `{"text":"` + msg.Text + `"`
	if attachment != "" {
		body += `,"attachments": [` + attachment + `]`
	}
	body += `}`

	return SendBody(msg.Url, []byte(body))
}

func SendBody(url string, bodyData []byte) error {
	body := bytes.NewReader(bodyData)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Err(err, "create slack post request")
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err, "create slack post response")
		return err
	}

	if resp.StatusCode != 200 {
		log.Error("status code != 200:", resp.StatusCode)
		return errors.New("invalid status code: " + strconv.Itoa(resp.StatusCode))
	}

	return nil
}

func GetAttachment(msg SlackMessage) string {
	if msg.Info == "" {
		return ""
	}

	blockText := "*info:* _" + msg.Info + "_\n"
	attachment := ""
	if msg.Color == "" {
		attachment = "{"
	} else {
		attachment = `{"color":"` + msg.Color + `",`
	}
	attachment += `"blocks": [{"type": "section","text": {"text": "` + blockText + `","type": "mrkdwn"}}]}`
	return attachment
}
