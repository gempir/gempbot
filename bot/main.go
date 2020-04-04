package main

import (
	"flag"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/api"
	"github.com/gempir/spamchamp/bot/collector"
	"github.com/gempir/spamchamp/bot/config"
	"github.com/gempir/spamchamp/bot/helix"
	"github.com/gempir/spamchamp/bot/stats"
)

var messageQueue = make(chan twitch.PrivateMessage)
var broadcastQueue = make(chan api.BroadcastMessage)

func main() {
	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	helixClient := helix.NewClient(cfg.ClientID)
	for _, id := range helixClient.GetTopChannels() {
		cfg.AddChannels(id)
	}

	bot := collector.NewBot(cfg, &helixClient, messageQueue)
	server := api.NewServer(cfg, &helixClient, broadcastQueue)
	broadcaster := stats.NewBroadcaster(messageQueue, broadcastQueue)



	go server.Start()
	go broadcaster.Start()

	bot.Connect()
}
