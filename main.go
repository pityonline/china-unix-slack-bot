package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: ./bot $SLACK_TOKEN\n")
		os.Exit(1)
	}

	ws, id := slackConnect(os.Args[1])
	fmt.Println("bot ready, ^C exits")

	for {
		m, err := getMessage(ws)
		if err != nil {
			log.Errorf("Get message Error: %#v", err)
		}

		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			parts := strings.Fields(m.Text)
			if len(parts) == 2 && parts[1] == "hi" {
				go func(m Message) {
					m.Text = greetting()
					postMessage(ws, m)
					log.Infof("%#v", m)
				}(m)
			} else if len(parts) == 3 && parts[1] == "ip" {
				api := "http://freeapi.ipip.net/"
				ip := parts[2]
				go func(m Message) {
					m.Text = ipQuery(api, ip)
					postMessage(ws, m)
					log.Infof("%#v", m)
				}(m)
			} else {
				m.Text = fmt.Sprintf("Sorry, it's not implemented yet.\n")
				postMessage(ws, m)
				log.Warnf("%#v", m)
			}
		}
	}
}

func greetting() string {
	return "Hello world!"
}

// ipQuery requests to an API with an IP address, return with location
func ipQuery(api, ip string) (loc string) {
	url := api + ip
	res, err := http.Get(url)

	if err != nil {
		log.Errorf("Request Error: %+v\n", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Errorf("Response Error: %+v\n", err)
	}

	loc = strings.TrimSpace(string(body))

	return fmt.Sprintf("IP: %v\nLocation: %v\nURL: %v\n", ip, loc, url)
}
