package emotechief

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/store"
)

type EmoteChief struct {
	cfg         *config.Config
	db          *store.Database
	helixClient *helix.Client
	httpClient  *http.Client
	chatClient  *chat.ChatClient
}

func NewEmoteChief(cfg *config.Config, db *store.Database, helixClient *helix.Client, chatClient *chat.ChatClient) *EmoteChief {
	return &EmoteChief{
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
		httpClient:  &http.Client{},
		chatClient:  chatClient,
	}
}
