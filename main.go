package main

import (
	"context"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	slackbot "github.com/flw-cn/go-slackbot"
	smartConfig "github.com/flw-cn/go-smart-config"
	"github.com/flw-cn/slack"
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

var Dumper = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
}

func main() {
	var config Config

	smartConfig.LoadConfig("Slack Bot", "v0.2.0", &config)

	api := slack.New(config.Token)
	bot, _ := slackbot.NewWithOpts(slackbot.WithClient(api))

	if config.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("Running in debug mode...")
		api.SetDebug(true)
	}

	bot.Hear("(?i)(hi|hello).*").MessageHandler(Hello)
	bot.Hear("(?i)(ping).*").MessageHandler(Ping)
	bot.Hear("(?i)(ip) .*").MessageHandler(QueryIP)

	bot.Run(false, nil)
}

func Hello(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	m := service.Greet()
	bot.Reply(evt, m, slackbot.WithTyping)
}

func Ping(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	m := service.Ping()
	bot.Reply(evt, m, slackbot.WithTyping)
}

func QueryIP(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	id := bot.BotUserID()
	id = "<@" + id + ">"
	text := strings.Replace(evt.Text, id, "", -1)
	strings.Trim(text, " ")
	parts := strings.Fields(text)

	api := "http://freeapi.ipip.net/"
	var m string
	if len(parts) != 2 {
		m = "Usage: ip <ip address>"
	} else {
		ip := parts[1]
		m = service.IPQuery(api, ip)
	}
	bot.Reply(evt, m, slackbot.WithTyping)
}
