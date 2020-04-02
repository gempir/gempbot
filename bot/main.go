package main

import (
	"flag"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/collector"
	"github.com/gempir/spamchamp/bot/config"
	"github.com/gempir/spamchamp/bot/helix"
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
	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)

	helixClient := helix.NewClient(cfg.ClientID)
	bot := collector.NewBot(cfg, &helixClient, messageQueue)

	go startStatsCollector()
	go startWebsocketServer()

	bot.Connect()
}
