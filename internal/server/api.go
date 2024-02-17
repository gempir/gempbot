package server

import (
	"github.com/gempir/gempbot/internal/auth"
	"github.com/gempir/gempbot/internal/bot"
	"github.com/gempir/gempbot/internal/channelpoint"
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/emotechief"
	"github.com/gempir/gempbot/internal/emoteservice"
	"github.com/gempir/gempbot/internal/eventsubmanager"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/store"
	"github.com/gempir/gempbot/internal/user"
	"github.com/gempir/gempbot/internal/ws"
	"github.com/gempir/gempbot/internal/ysweet"
)

type Api struct {
	db                  *store.Database
	cfg                 *config.Config
	helixClient         helixclient.Client
	userAdmin           *user.UserAdmin
	authClient          *auth.Auth
	bot                 *bot.Bot
	emoteChief          *emotechief.EmoteChief
	eventsubManager     *eventsubmanager.EventsubManager
	channelPointManager *channelpoint.ChannelPointManager
	sevenTvClient       emoteservice.ApiClient
	wsHandler           *ws.WsHandler
	tokenFactory        *ysweet.Factory
}

func NewApi(cfg *config.Config, db *store.Database, helixClient helixclient.Client, userAdmin *user.UserAdmin, authClient *auth.Auth, bot *bot.Bot, emoteChief *emotechief.EmoteChief, eventsubManager *eventsubmanager.EventsubManager, channelPointManager *channelpoint.ChannelPointManager, sevenTvClient emoteservice.ApiClient, wsHandler *ws.WsHandler, tokenFactory *ysweet.Factory) *Api {
	return &Api{
		db:                  db,
		cfg:                 cfg,
		helixClient:         helixClient,
		userAdmin:           userAdmin,
		authClient:          authClient,
		bot:                 bot,
		emoteChief:          emoteChief,
		eventsubManager:     eventsubManager,
		channelPointManager: channelPointManager,
		sevenTvClient:       sevenTvClient,
		wsHandler:           wsHandler,
		tokenFactory:        tokenFactory,
	}
}
