package helix

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	nickHelix "github.com/nicklaw5/helix"
)

// Client wrapper for helix
type Client struct {
	clientID       string
	clientSecret   string
	eventSubSecret string
	Client         *nickHelix.Client
	AppAccessToken nickHelix.AccessCredentials
	db             *store.Database
	httpClient     *http.Client
}

var (
	userCacheByID       map[string]*UserData
	userCacheByUsername map[string]*UserData
)

func init() {
	userCacheByID = map[string]*UserData{}
	userCacheByUsername = map[string]*UserData{}
}

var scopes = []string{"channel:read:redemptions", "channel:manage:redemptions", "channel:read:predictions", "channel:manage:predictions"}

// NewClient Create helix client
func NewClient(cfg *config.Config, db *store.Database) *Client {
	client, err := nickHelix.NewClient(&nickHelix.Options{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURI:  cfg.ApiBaseUrl + "/api/callback",
	})
	if err != nil {
		panic(err)
	}

	token, err := db.GetAppAccessToken()
	if err != nil || time.Since(token.UpdatedAt) > 24*time.Hour {
		log.Info("App AccessToken not found or older than 24hours")
		resp, err := client.RequestAppAccessToken(scopes)
		if err != nil {
			panic(err)
		}
		log.Infof("Requested access token, response: %d, expires in: %d", resp.StatusCode, resp.Data.ExpiresIn)
		client.SetAppAccessToken(resp.Data.AccessToken)
		err = db.SaveAppAccessToken(context.Background(), resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "), resp.Data.ExpiresIn)
		if err != nil {
			log.Errorf("Failure saving app access token: %s", err.Error())
		}
		token = store.AppAccessToken{AccessToken: resp.Data.AccessToken, RefreshToken: resp.Data.RefreshToken, Scopes: strings.Join(resp.Data.Scopes, " "), ExpiresIn: resp.Data.ExpiresIn}
	} else {
		client.SetAppAccessToken(token.AccessToken)
	}

	return &Client{
		clientID:       cfg.ClientID,
		clientSecret:   cfg.ClientSecret,
		eventSubSecret: cfg.Secret,
		Client:         client,
		AppAccessToken: nickHelix.AccessCredentials{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken, Scopes: strings.Split(token.Scopes, " "), ExpiresIn: token.ExpiresIn},
		db:             db,
		httpClient:     &http.Client{},
	}
}

// StartRefreshTokenRoutine refresh our token
func (c *Client) StartRefreshTokenRoutine() {
	ticker := time.NewTicker(24 * time.Hour)

	for range ticker.C {
		resp, err := c.Client.RequestAppAccessToken(scopes)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Requested access token from routine, response: %d, expires in: %d", resp.StatusCode, resp.Data.ExpiresIn)

		c.SetAppAccessToken(context.Background(), resp.Data)
	}
}

func (c *Client) SetAppAccessToken(ctx context.Context, token nickHelix.AccessCredentials) {
	c.AppAccessToken = token
	c.Client.SetAppAccessToken(token.AccessToken)
	err := c.db.SaveAppAccessToken(ctx, token.AccessToken, token.RefreshToken, strings.Join(token.Scopes, " "), token.ExpiresIn)
	if err != nil {
		log.Errorf("Failure saving app access token: %s", err.Error())
	}
}
