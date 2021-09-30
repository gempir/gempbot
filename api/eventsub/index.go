package eventsub

import (
	"fmt"
	"net/http"

	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/log"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	cfg := config.FromEnv()
	eventSubManager := eventsub.NewEventSubManager(cfg)

	body, err := eventSubManager.HandleWebhook(w, r)
	if err != nil {
		http.Error(w, err.Error(), err.Status())
		return
	}
	if len(body) == 0 {
		return
	}

	log.Info(body)

	fmt.Fprintf(w, "Hello")
}
