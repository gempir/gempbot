package helix

import (
	"net/http"
	"time"

	"github.com/gempir/bitraft/pkg/log"
	helixClient "github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Client wrapper for helix
type Client struct {
	clientID               string
	clientSecret           string
	eventSubSecret         string
	Client                 *helixClient.Client
	httpClient             *http.Client
	helixApiResponseStatus *prometheus.CounterVec
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
func NewClient(clientID, clientSecret, redirectURI string, eventSubSecret string) *Client {
	client, err := helixClient.NewClient(&helixClient.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
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
		clientID:       clientID,
		clientSecret:   clientSecret,
		eventSubSecret: eventSubSecret,
		Client:         client,
		httpClient:     &http.Client{},
		helixApiResponseStatus: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "helix_responses",
				Help: "The status codes of all twitch api responses",
			}, []string{
				"code",
				"function",
			},
		),
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
