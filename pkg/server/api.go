package server

import (
	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/bot"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/user"
)

type Api struct {
	db          *store.Database
	cfg         *config.Config
	helixClient *helixclient.Client
	userAdmin   *user.UserAdmin
	authClient  *auth.Auth
	bot         *bot.Bot
}

func NewApi(cfg *config.Config, db *store.Database, helixClient *helixclient.Client, userAdmin *user.UserAdmin, authClient *auth.Auth, bot *bot.Bot) *Api {
	return &Api{
		db:          db,
		cfg:         cfg,
		helixClient: helixClient,
		userAdmin:   userAdmin,
		authClient:  authClient,
		bot:         bot,
	}
}
