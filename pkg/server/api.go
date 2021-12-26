package server

import (
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/bot"
	"github.com/gempir/gempbot/pkg/channelpoint"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/emotechief"
	"github.com/gempir/gempbot/pkg/eventsub"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
)

type Api struct {
	db                          *store.Database
	cfg                         *config.Config
	helixClient                 *helixclient.Client
	userAdmin                   *user.UserAdmin
	authClient                  *auth.Auth
	bot                         *bot.Bot
	emoteChief                  *emotechief.EmoteChief
	eventsubManager             *eventsub.EventsubManager
	eventsubSubscriptionManager *eventsub.SubscriptionManager
	channelPointManager         *channelpoint.ChannelPointManager
}

func NewApi(cfg *config.Config, db *store.Database, helixClient *helixclient.Client, userAdmin *user.UserAdmin, authClient *auth.Auth, bot *bot.Bot, emoteChief *emotechief.EmoteChief, eventsubManager *eventsub.EventsubManager, eventsubSubscriptionManager *eventsub.SubscriptionManager, channelPointManager *channelpoint.ChannelPointManager) *Api {
	return &Api{
		db:                          db,
		cfg:                         cfg,
		helixClient:                 helixClient,
		userAdmin:                   userAdmin,
		authClient:                  authClient,
		bot:                         bot,
		emoteChief:                  emoteChief,
		eventsubManager:             eventsubManager,
		eventsubSubscriptionManager: eventsubSubscriptionManager,
		channelPointManager:         channelPointManager,
	}
}
