package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gempir/bitraft/services/api/emotechief"
	echoPrometheus "github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/slice"
	"github.com/gempir/bitraft/pkg/store"
	nickHelix "github.com/nicklaw5/helix"
)

// Server api server
type Server struct {
	cfg         *config.Config
	helixClient *helix.Client
	store       *store.Redis
	db          *store.Database
	emotechief  *emotechief.EmoteChief
}

// NewServer create api Server
func NewServer(cfg *config.Config, helixClient *helix.Client, store *store.Redis, db *store.Database, emotechief *emotechief.EmoteChief) Server {
	return Server{
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
		store:       store,
		emotechief:  emotechief,
	}
}

func (s *Server) Start() {
	go s.syncSubscriptions()
	go s.tokenRefreshRoutine()

	e := echo.New()
	e.HideBanner = true
	p := echoPrometheus.NewPrometheus("api", nil)
	p.Use(e)

	e.GET("/api/callback", s.handleCallback)
	e.POST("/api/redemption", s.handleChannelPointsRedemption)
	e.POST("/api/prediction/begin", s.handlePredictionBegin)
	e.POST("/api/prediction/lock", s.handlePredictionLock)
	e.POST("/api/prediction/end", s.handlePredictionEnd)
	e.GET("/api/userConfig", s.handleUserConfig)
	e.POST("/api/userConfig", s.handleUserConfig)

	e.GET("/api/reward/:userID", s.handleRewardRead)
	e.GET("/api/reward/:userID/type/:type", s.handleRewardSingleRead)

	e.DELETE("/api/reward/:userID/type/:type", s.handleRewardDeletion)

	e.POST("/api/reward/:userID", s.handleRewardCreateOrUpdate)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{s.cfg.WebBaseUrl},
	}))
	e.Logger.Fatal(e.Start(":8035"))
}

func (s *Server) syncSubscriptions() {
	resp, err := s.helixClient.Client.GetEventSubSubscriptions(&nickHelix.EventSubSubscriptionsParams{})
	if err != nil {
		log.Errorf("Failed to get subscriptions: %s", err)
		return
	}

	log.Infof("Found %d/%d total cost subscriptions, syncing to DB", resp.Data.TotalCost, resp.Data.MaxTotalCost)
	subscribed := []string{}

	for _, sub := range resp.Data.EventSubSubscriptions {
		if !strings.Contains(sub.Transport.Callback, s.cfg.WebhookApiBaseUrl) || sub.Status == nickHelix.EventSubStatusFailed {
			err := s.removeEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Type, "bad EventSub subscription, unsubscribing")
			if err != nil {
				log.Errorf("Failed to unsubscribe %s error: %s", sub.Condition.BroadcasterUserID, err.Error())
			}
			continue
		}

		_, err = s.db.GetEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Type)
		if err != nil {
			log.Infof("Found unknown subscription, adding %s", sub.Condition.BroadcasterUserID)
			s.db.AddEventSubSubscription(sub.Condition.BroadcasterUserID, sub.ID, sub.Version, sub.Type)
		}
		subscribed = append(subscribed, sub.Condition.BroadcasterUserID+sub.Type)
	}

	rewards := s.db.GetDistinctRewardsPerUser()
	log.Infof("Found %d total distinct rewards, checking missing subscriptions", len(rewards))

	for _, dbReward := range rewards {
		if !slice.Contains(subscribed, dbReward.OwnerTwitchID+nickHelix.EventSubTypeChannelPointsCustomRewardRedemptionAdd) {
			err := s.subscribeChannelPoints(dbReward.OwnerTwitchID)
			if err != nil {
				log.Infof("Removing reward for user %s because we didn't get permission to subscribe eventsub", dbReward.OwnerTwitchID)
				s.db.DeleteChannelPointReward(dbReward.OwnerTwitchID, dbReward.Type)
			}
		}
	}
}

func (s *Server) handleWebhook(c echo.Context, response interface{}) (bool, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Error(err)
		return true, echo.NewHTTPError(http.StatusBadRequest, "Failed reading body")
	}

	verified := nickHelix.VerifyEventSubNotification(s.cfg.Secret, c.Request().Header, string(body))
	if !verified {
		log.Errorf("Failed verification %s", c.Request().Header.Get("Twitch-Eventsub-Message-Id"))
		return true, echo.NewHTTPError(http.StatusPreconditionFailed, "failed verfication")
	}

	if c.Request().Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		return true, s.handleChallenge(c, body)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return true, echo.NewHTTPError(http.StatusBadRequest, "Failed decoding body: "+err.Error())
	}

	return false, nil
}

func (s *Server) handleChallenge(c echo.Context, body []byte) error {
	var event struct {
		Challenge string `json:"challenge"`
	}
	err := json.Unmarshal(body, &event)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to handle challenge: "+err.Error(), http.StatusBadRequest))
	}

	log.Infof("Challenge success: %s", event.Challenge)
	return c.String(http.StatusOK, event.Challenge)
}
