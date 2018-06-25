package main

import (
	"fmt"
	"os"
	"strings"

	smartConfig "github.com/flw-cn/go-smart-config"
	"github.com/pityonline/china-unix-slack-bot/service"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Debug bool   `flag:"d|false|debug mode, default to 'false'"`
	Token string `flag:"t||must provide your {SLACK_TOKEN} here"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var config Config

	smartConfig.LoadConfig("Slack Bot", "v0.2.0", &config)

	if config.Debug {
		log.Println("Running in debug mode...")
		log.SetLevel(log.DebugLevel)
	}

	runBot(config.Token)
}

// runBot runs the bot service
func runBot(token string) error {
	ws, id := slackConnect(token)
	log.Println("Bot ready, ^C exits")

	for {
		m, err := getMessage(ws)
		if err != nil {
			log.Errorf("Get message Error: %#v", err)
		}

		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			parts := strings.Fields(m.Text)
			if len(parts) == 2 && parts[1] == "hi" {
				go func(m Message) {
					m.Text = service.Greet()
					postMessage(ws, m)
					log.Infof("%#v", m)
				}(m)
			} else if len(parts) == 2 && parts[1] == "ping" {
				go func(m Message) {
					m.Text = service.Ping()
					postMessage(ws, m)
					log.Infof("%#v", m)
				}(m)
			} else if len(parts) == 3 && parts[1] == "ip" {
				api := "http://freeapi.ipip.net/"
				ip := parts[2]
				go func(m Message) {
					m.Text = service.IPQuery(api, ip)
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
