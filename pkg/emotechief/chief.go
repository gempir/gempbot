package emotechief

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/store"
)

type EmoteChief struct {
	cfg         *config.Config
	db          *store.Database
	helixClient *helixclient.Client
	httpClient  *http.Client
	chatClient  *chat.ChatClient
}

func NewEmoteChief(cfg *config.Config, db *store.Database, helixClient *helixclient.Client, chatClient *chat.ChatClient) *EmoteChief {
	return &EmoteChief{
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
		httpClient:  &http.Client{},
		chatClient:  chatClient,
	}
}
