package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var debug bool
	var token string

	// app information
	app := cli.NewApp()
	app.Name = "China Unix Slack Bot"
	app.Version = "v0.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "pityonline",
			Email: "pityonline@gmail.com",
		},
	}
	app.Usage = "A bot service for ChinaUnix slack team"

	// global flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "debug mode",
			Destination: &debug,
		},
		cli.StringFlag{
			Name: "token, t",
			// FIXME: EnvVar not taken if not passed, even it's in env
			EnvVar:      "SLACK_TOKEN",
			Usage:       "must provide your slack token",
			Destination: &token,
		},
	}

	// FIXME: there's must be an action before app.Run, or it prints help
	app.Action = func(ctx *cli.Context) error {
		fmt.Printf("%s version %s\n", app.Name, app.Version)
		return nil
	}

	// run app
	err := app.Run(os.Args)

	if debug == true {
		fmt.Println("Running in debug mode...")
		log.SetLevel(log.DebugLevel)
	}

	if err != nil {
		log.Fatalf("%#v", err)
	}

	app.Action = runBot(token)
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
