package emotechief

import (
	"github.com/gempir/gempbot/internal/chat"
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/emoteservice"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/store"
)

type EmoteChief struct {
	cfg           *config.Config
	db            store.Store
	helixClient   helixclient.Client
	chatClient    *chat.ChatClient
	sevenTvClient emoteservice.ApiClient
}

func NewEmoteChief(cfg *config.Config, db store.Store, helixClient helixclient.Client, chatClient *chat.ChatClient, sevenTvClient emoteservice.ApiClient) *EmoteChief {
	return &EmoteChief{
		cfg:           cfg,
		db:            db,
		helixClient:   helixClient,
		chatClient:    chatClient,
		sevenTvClient: sevenTvClient,
	}
}
