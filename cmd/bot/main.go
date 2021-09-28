package main

import (
	"fmt"
	"net/http"

	"github.com/gempir/gempbot/cmd/bot/collector"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/store"
)

func main() {
	cfg := config.FromEnv()
	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg, db)
	go helixClient.StartRefreshTokenRoutine()

	bot := collector.NewBot(cfg, db, helixClient)
	go bot.Connect()

	http.HandleFunc("/", status)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()

	<-bot.Done
}

func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "ok\n")
}
