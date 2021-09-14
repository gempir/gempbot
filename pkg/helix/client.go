package helix

import (
	"net/http"
	"time"

	"github.com/gempir/bot/pkg/config"
	"github.com/gempir/bot/pkg/log"
	helixClient "github.com/nicklaw5/helix"
)

// Client wrapper for helix
type Client struct {
	clientID       string
	clientSecret   string
	eventSubSecret string
	Client         *helixClient.Client
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

// NewClient Create helix client
func NewClient(cfg *config.Config) *Client {
	client, err := helixClient.NewClient(&helixClient.Options{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURI:  cfg.ApiBaseUrl + "/api/callback",
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.RequestAppAccessToken([]string{"channel:read:redemptions"})
	if err != nil {
		panic(err)
	}
	log.Infof("Requested access token, response: %d, expires in: %d", resp.StatusCode, resp.Data.ExpiresIn)
	client.SetAppAccessToken(resp.Data.AccessToken)

	return &Client{
		clientID:       cfg.ClientID,
		clientSecret:   cfg.ClientSecret,
		eventSubSecret: cfg.Secret,
		Client:         client,
		httpClient:     &http.Client{},
	}
}

// StartRefreshTokenRoutine refresh our token
func (c *Client) StartRefreshTokenRoutine() {
	ticker := time.NewTicker(24 * time.Hour)

	for range ticker.C {
		resp, err := c.Client.RequestAppAccessToken([]string{})
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Requested access token from routine, response: %d, expires in: %d", resp.StatusCode, resp.Data.ExpiresIn)

		c.Client.SetAppAccessToken(resp.Data.AccessToken)
	}
}
