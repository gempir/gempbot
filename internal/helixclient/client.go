package helixclient

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/jellydator/ttlcache/v2"
	"github.com/nicklaw5/helix/v2"
)

type Client interface {
	StartRefreshTokenRoutine()
	RefreshToken(token store.UserAccessToken) error
	CreateEventSubSubscription(userID string, webHookUrl string, subType string) (*helix.EventSubSubscriptionsResponse, error)
	CreateRewardEventSubSubscription(userID, webHookUrl, subType, rewardID string, false bool) (*helix.EventSubSubscriptionsResponse, error)
	RemoveEventSubSubscription(id string) (*helix.RemoveEventSubSubscriptionParamsResponse, error)
	GetEventSubSubscriptions(params *helix.EventSubSubscriptionsParams) (*helix.EventSubSubscriptionsResponse, error)
	GetAllSubscriptions(eventType string) []helix.EventSubSubscription
	GetPredictions(params *helix.PredictionsParams) (*helix.PredictionsResponse, error)
	EndPrediction(params *helix.EndPredictionParams) (*helix.PredictionsResponse, error)
	CreatePrediction(params *helix.CreatePredictionParams) (*helix.PredictionsResponse, error)
	CreateOrUpdateReward(userID string, reward CreateCustomRewardRequest, rewardID string) (*helix.ChannelCustomReward, error)
	UpdateRedemptionStatus(broadcasterID, rewardID string, redemptionID string, statusSuccess bool) error
	DeleteReward(userID string, rewardID string) error
	GetUsersByUserIds(userIDs []string) (map[string]UserData, error)
	GetUsersByUsernames(usernames []string) (map[string]UserData, error)
	GetUserByUsername(username string) (UserData, error)
	GetUserByUserID(userID string) (UserData, error)
	SetUserAccessToken(token string)
	ValidateToken(accessToken string) (bool, *helix.ValidateTokenResponse, error)
	RequestUserAccessToken(code string) (*helix.UserAccessTokenResponse, error)
	SendChatMessage(params *helix.SendChatMessageParams) (*SendChatMessageResponse, error)
}

// Client wrapper for helix
type HelixClient struct {
	clientID        string
	clientSecret    string
	eventSubSecret  string
	Client          *helix.Client
	AppAccessToken  helix.AccessCredentials
	db              store.Store
	httpClient      *http.Client
	refreshTtlCache *ttlcache.Cache
}

var (
	cacheMutex          *sync.Mutex
	userCacheByID       map[string]*UserData
	userCacheByUsername map[string]*UserData
)

func init() {
	cacheMutex = &sync.Mutex{}
	userCacheByID = map[string]*UserData{}
	userCacheByUsername = map[string]*UserData{}
}

const TWITCH_API = "https://api.twitch.tv/"

var scopes = []string{"channel:read:redemptions", "channel:manage:redemptions", "channel:read:predictions", "channel:manage:predictions moderation:read channel:bot user:write:chat user:bot"}

// NewClient Create helix client
func NewClient(cfg *config.Config, db store.Store) *HelixClient {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURI:  cfg.WebBaseUrl + "/api/callback",
	})
	if err != nil {
		panic(err)
	}
	cache := ttlcache.NewCache()
	err = cache.SetTTL(time.Second * 30)
	if err != nil {
		panic(err)
	}
	token := setOrUpdateAccessToken(client, db)

	return &HelixClient{
		clientID:        cfg.ClientID,
		clientSecret:    cfg.ClientSecret,
		eventSubSecret:  cfg.Secret,
		Client:          client,
		AppAccessToken:  helix.AccessCredentials{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken, Scopes: strings.Split(token.Scopes, " "), ExpiresIn: token.ExpiresIn},
		db:              db,
		httpClient:      &http.Client{},
		refreshTtlCache: cache,
	}
}

// StartRefreshTokenRoutine refresh our token
func (c *HelixClient) StartRefreshTokenRoutine() {
	setOrUpdateAccessToken(c.Client, c.db)
	go func() {
		for range time.NewTicker(1 * time.Hour).C {
			setOrUpdateAccessToken(c.Client, c.db)
		}
	}()

	c.refreshUserAccessTokens()
	go func() {
		for range time.NewTicker(1 * time.Hour).C {
			c.refreshUserAccessTokens()
		}
	}()
}

func (c *HelixClient) refreshUserAccessTokens() {
	tokens := c.db.GetAllUserAccessToken()
	for _, token := range tokens {
		if time.Since(token.UpdatedAt) > 3*time.Hour {
			err := c.RefreshToken(token)
			if err != nil {
				log.Errorf("failed to refresh token for user %s %s", token.OwnerTwitchID, err)
			} else {
				log.Infof("refreshed token for user %s", token.OwnerTwitchID)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func (c *HelixClient) refreshUserAccessToken(userID string) error {
	if item, _ := c.refreshTtlCache.Get(userID); item != nil {
		return fmt.Errorf("already refreshing token for user %s in last 30 seconds", userID)
	}
	c.refreshTtlCache.Set(userID, true)

	token, err := c.db.GetUserAccessToken(userID)
	if err != nil {
		return err
	}
	err = c.RefreshToken(token)
	if err != nil {
		return err
	} else {
		log.Infof("refreshed token for user %s", token.OwnerTwitchID)
	}

	return nil
}

func (c *HelixClient) SetAppAccessToken(ctx context.Context, token helix.AccessCredentials) {
	c.AppAccessToken = token
	c.Client.SetAppAccessToken(token.AccessToken)
	err := c.db.SaveAppAccessToken(ctx, token.AccessToken, token.RefreshToken, strings.Join(token.Scopes, " "), token.ExpiresIn)
	if err != nil {
		log.Errorf("Failure saving app access token: %s", err.Error())
	}
}

func setOrUpdateAccessToken(client *helix.Client, db store.Store) store.AppAccessToken {
	token, err := db.GetAppAccessToken()
	if err != nil || time.Since(token.UpdatedAt) > 24*time.Hour {
		log.Info("App AccessToken not found or older than 24 hours")
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

	return token
}

func (c *HelixClient) RefreshToken(token store.UserAccessToken) error {
	resp, err := c.Client.RefreshUserAccessToken(token.RefreshToken)
	if err != nil {
		return err
	}

	err = c.db.SaveUserAccessToken(context.Background(), token.OwnerTwitchID, resp.Data.AccessToken, resp.Data.RefreshToken, strings.Join(resp.Data.Scopes, " "))
	if err != nil {
		return err
	}

	return nil
}

func (c *HelixClient) SetUserAccessToken(token string) {
	c.Client.SetUserAccessToken(token)
}

func (c *HelixClient) ValidateToken(accessToken string) (bool, *helix.ValidateTokenResponse, error) {
	return c.Client.ValidateToken(accessToken)
}

func (c *HelixClient) RequestUserAccessToken(code string) (*helix.UserAccessTokenResponse, error) {
	return c.Client.RequestUserAccessToken(code)
}
