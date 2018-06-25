package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	smartConfig "github.com/flw-cn/go-smart-config"
	"github.com/nlopes/slack"
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

var selfID string

var Dumper = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
}

func main() {
	var config Config

	smartConfig.LoadConfig("Slack Bot", "v0.2.0", &config)

	client := slack.New(config.Token)

	if config.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Running in debug mode...")
		client.SetDebug(true)
	}

	rtm := client.NewRTM()
	go rtm.ManageConnection()

	mainLoop(rtm)
}

func mainLoop(rtm *slack.RTM) {
	for m := range rtm.IncomingEvents {
		switch ev := m.Data.(type) {
		case *slack.ConnectingEvent:
			log.Info("Connecting...")

		case *slack.ConnectedEvent:
			log.Infof("Connected to %s<%s>, as known as %s<%s>",
				ev.Info.Team.Name, ev.Info.Team.ID,
				ev.Info.User.Name, ev.Info.User.ID,
			)

			selfID = ev.Info.User.ID

		case *slack.HelloEvent:
			// Ignore hello
			log.Info("Received hello message which come from Slack. Online now.")

		case *slack.MessageEvent:
			log.Info("Message:", ev.Text)
			log.Debug("Message:", Dumper.Sdump(ev))
			runBot(rtm, ev)

		case *slack.RTMError:
			log.Debug("Error:", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Debug("Invalid credentials:", ev)

		case *slack.LatencyReport:
			log.Debug("LatencyReport:", ev.Value)

		default:
			// Ignore other events..
			log.Debug("Unexpected:", ev)
		}
	}
}

func runBot(rtm *slack.RTM, ev *slack.MessageEvent) {
	text := ev.Text
	id := "<@" + selfID + ">"
	if !strings.Contains(text, id) {
		return
	}

	text = strings.Replace(text, id, "", -1)
	strings.Trim(text, " ")

	parts := strings.Fields(text)
	if len(parts) == 1 && parts[0] == "hi" {
		go func() {
			m := service.Greet()
			rtm.SendMessage(rtm.NewOutgoingMessage(m, ev.Channel))
			log.Infof("%#v", m)
		}()
	} else if len(parts) == 1 && parts[0] == "ping" {
		go func() {
			m := service.Ping()
			rtm.SendMessage(rtm.NewOutgoingMessage(m, ev.Channel))
			log.Infof("%#v", m)
		}()
	} else if len(parts) == 2 && parts[0] == "ip" {
		api := "http://freeapi.ipip.net/"
		ip := parts[1]
		go func() {
			m := service.IPQuery(api, ip)
			rtm.SendMessage(rtm.NewOutgoingMessage(m, ev.Channel))
			log.Infof("%#v", m)
		}()
	} else {
		go func() {
			m := fmt.Sprintf("Sorry, it's not implemented yet.\n")
			rtm.SendMessage(rtm.NewOutgoingMessage(m, ev.Channel))
			log.Warnf("%#v", m)
		}()
	}
}
