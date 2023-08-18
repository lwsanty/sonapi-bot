package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/itchyny/json2yaml"
	tb "gopkg.in/telebot.v3"
)

const (
	endpoint = "https://api.sonapi.ee/v1/"

	credits = `Bot uses [Sõnad API](https://www.sonapi.ee)([code](https://github.com/BenediktGeiger/sonad-api)) built by [BenediktGeiger](https://github.com/BenediktGeiger).
Contribute https://github.com/lwsanty/sonapi-bot`
	usageText = `/help for commands list
[word]			- infos about word
ps [word]		- part of speech of word
wf [word]		- word forms of word
ms [word]		- meanings of word`
)

var (
	usage = fmt.Sprintf("```\n%s\n```", usageText)
	help  = usage + "\n\n" + credits
)

var methods = map[string]string{
	"ps": "partofspeech",
	"wf": "wordforms",
	"ms": "meanings",
}

func main() {
	pref := tb.Settings{
		Token:  os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tb.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	b.Handle("/start", func(c tb.Context) error {
		return c.Reply(help, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})
	b.Handle("/help", func(c tb.Context) error {
		return c.Reply(help, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})
	b.Handle(tb.OnText, func(c tb.Context) error {
		// TODO: message validation
		var method string
		text := c.Text()
		for k, m := range methods {
			prefix := k + " "
			if !strings.HasPrefix(text, prefix) {
				continue
			}
			method = m
			text = strings.TrimPrefix(text, prefix)
			break
		}

		respText, err := doReq(text, method)
		if err != nil {
			log.Fatal(err)
			return err
		}

		return c.Reply(respText, &tb.SendOptions{ParseMode: tb.ModeMarkdown})
	})
	b.Start()
}

func doReq(text, method string) (string, error) {
	fEndpoint := endpoint + text
	if method != "" {
		fEndpoint += "/" + method
	}
	u, err := url.Parse(fEndpoint)
	if err != nil {
		return "", err
	}
	u.RawQuery = u.Query().Encode()
	// TODO: removeme
	log.Print(u.String())

	res, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	// it's ok to read the whole response body in this case
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return jsonToYaml(string(data)), nil
}

func jsonToYaml(s string) string {
	input := strings.NewReader(s)
	var output strings.Builder
	if err := json2yaml.Convert(&output, input); err != nil {
		log.Fatalln(err)
	}

	return fmt.Sprintf("```\n%s\n```", output.String())
}
