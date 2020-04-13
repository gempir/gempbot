package main

import (
	"flag"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/api"
	"github.com/gempir/spamchamp/bot/collector"
	"github.com/gempir/spamchamp/bot/config"
	"github.com/gempir/spamchamp/bot/helix"
	"github.com/gempir/spamchamp/bot/stats"
	"github.com/gempir/spamchamp/bot/store"
)

var messageQueue = make(chan twitch.PrivateMessage)
var broadcastQueue = make(chan api.BroadcastMessage)

func main() {
	configFile := flag.String("config", "config.json", "json config file")
	flag.Parse()

	cfg := config.NewConfig(*configFile)
	rStore := store.NewStore()
	helixClient := helix.NewClient(cfg.ClientID)

	bot := collector.NewBot(cfg, &helixClient, rStore, messageQueue)
	bot.LoadTopChannelsAndJoin()
	server := api.NewServer(cfg, &helixClient, broadcastQueue)
	broadcaster := stats.NewBroadcaster(messageQueue, broadcastQueue, rStore, bot)

	go server.Start()
	go broadcaster.Start()

	bot.Connect()
}
