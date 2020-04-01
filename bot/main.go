package main

import (
	twitch "github.com/gempir/go-twitch-irc/v2"
)

var messageQueue = make(chan twitch.PrivateMessage)

type socketMessage struct {
	Channels map[string]frontendStats `json:"channels"`
}

type frontendStats struct {
	ChannelName       string `json:"channelName"`
	MessagesPerSecond int    `json:"messagesPerSecond"`
}

func main() {
	client := twitch.NewClient("justinfan123123", "oauth:123123123")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		messageQueue <- message
	})

	client.Join("tmiloadtesting2")
	client.Join("xqcow")
	client.Join("lirik")
	client.Join("drdisrepect")
	client.Join("esl_csgo")
	client.Join("loltyler1")
	client.Join("forsen")

	go client.Connect()

	go startWebsocketServer()
	startStatsCollector()
}
