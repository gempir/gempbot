package main

import (
	"context"
	"strings"
	"time"

	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
)

var (
	cfg         *config.Config
	db          *store.Database
	helixClient *helixclient.Client
)

func main() {
	cfg = config.FromEnv()
	db = store.NewDatabase(cfg)
	helixClient = helixclient.NewClient(cfg, db)

	tokens := db.GetAllUserAccessToken()

	for _, token := range tokens {
		err := refreshToken(token)
		if err != nil {
			log.Errorf("failed to refresh token for user %s %s", token.OwnerTwitchID, err)
		} else {
			log.Infof("refreshed token for user %s", token.OwnerTwitchID)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func refreshToken(token store.UserAccessToken) error {
	resp, err := helixClient.Client.RefreshUserAccessToken(token.RefreshToken)
	if err != nil {
		return err
	}

	err = db.SaveUserAccessToken(context.Background(), token.OwnerTwitchID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
	if err != nil {
		return err
	}

	return nil
}
