package main

import (
	"flag"

	"github.com/nice-pink/goutil/pkg/notify"
)

func main() {
	url := flag.String("url", "", "Slack webhook url")
	flag.Parse()

	msg := notify.SlackMessage{
		Url:   *url,
		Text:  "this is the text",
		Info:  "this is the info",
		Color: "#334455",
	}
	notify.Send(msg)
}
