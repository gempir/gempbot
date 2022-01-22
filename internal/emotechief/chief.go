package emotechief

import (
	"github.com/gempir/gempbot/internal/chat"
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/store"
)

type EmoteChief struct {
	cfg         *config.Config
	db          *store.Database
	helixClient *helixclient.Client
	chatClient  *chat.ChatClient
}

func NewEmoteChief(cfg *config.Config, db *store.Database, helixClient *helixclient.Client, chatClient *chat.ChatClient) *EmoteChief {
	return &EmoteChief{
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
		chatClient:  chatClient,
	}
}
