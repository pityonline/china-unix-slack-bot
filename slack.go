package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type responseRtmStart struct {
	Ok    bool         `json:"ok"`
	Error string       `json:"error"`
	URL   string       `json:"url"`
	Self  responseSelf `json:"self"`
}

type responseSelf struct {
	ID string `json:"id"`
}

func slackStart(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Slack RTM Start Error: %#v", err)
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Errorf("Response Body Error: %#v", err)
	}
	var respObj responseRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		log.Errorf("JSON Unmarshal Error: %#v", err)
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.URL
	id = respObj.Self.ID
	return
}

// Message represents a message event
//
// API document: https://api.slack.com/events/message
type Message struct {
	ID      uint64 `json:"id"`
	User    string `json:"user"`
	Type    string `json:"type"`
	SubType string `json:"subtype"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	TS      string `json:"ts"`
}

func getMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

var counter uint64

func postMessage(ws *websocket.Conn, m Message) error {
	m.ID = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, m)
}

func slackConnect(token string) (*websocket.Conn, string) {
	wsurl, id, err := slackStart(token)
	log.Debugf("Slack token: %+v", token)
	if err != nil {
		log.Fatalf("Slack Start Error: %#v", err)
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		log.Fatalf("Slack Dail Error: %#v", err)
	}

	return ws, id
}
