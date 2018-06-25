# China Unix Slack Bot

[![Build Status](https://travis-ci.com/pityonline/china-unix-slack-bot.svg?branch=master)](https://travis-ci.com/pityonline/china-unix-slack-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/pityonline/china-unix-slack-bot)](https://goreportcard.com/report/github.com/pityonline/china-unix-slack-bot)

This bot is written in Go, inspired by [rapidloop/mybot](https://github.com/rapidloop/mybot).

## Usage

If you have Go installed, use `go get github.com/pityonline/china-unix-slack-bot` to install this app.

Specify the Slack token to run it (or make the token as a exported variable):

`china-unix-slack-bot -t $SLACK_TOKEN`

It will start the bot with some logs, use ^C to exit.
