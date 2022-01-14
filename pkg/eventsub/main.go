package eventsub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/emotechief"
	"github.com/gempir/gempbot/pkg/helixclient"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/nicklaw5/helix/v2"
)

type EventsubManager struct {
	cfg         *config.Config
	helixClient *helixclient.Client
	db          *store.Database
	emoteChief  *emotechief.EmoteChief
	chatClient  *chat.ChatClient
	ttlCache    *ttlcache.Cache
}

func NewEventsubManager(cfg *config.Config, helixClient *helixclient.Client, db *store.Database, emoteChief *emotechief.EmoteChief, bot *chat.ChatClient) *EventsubManager {
	cache := ttlcache.NewCache()
	err := cache.SetTTL(time.Second * 60)
	if err != nil {
		panic(err)
	}

	return &EventsubManager{
		cfg:         cfg,
		helixClient: helixClient,
		db:          db,
		emoteChief:  emoteChief,
		chatClient:  bot,
		ttlCache:    cache,
	}
}

type eventSubNotification struct {
	Subscription helix.EventSubSubscription `json:"subscription"`
	Challenge    string                     `json:"challenge"`
	Event        json.RawMessage            `json:"event"`
}

func (esm *EventsubManager) HandleWebhook(w http.ResponseWriter, r *http.Request) (event []byte, apiErr api.Error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		return []byte{}, api.NewApiError(http.StatusBadRequest, err)
	}

	verified := helix.VerifyEventSubNotification(esm.cfg.Secret, r.Header, string(body))
	if !verified {
		log.Errorf("Failed verification %s", r.Header.Get("Twitch-Eventsub-Message-Id"))
		return []byte{}, api.NewApiError(http.StatusPreconditionFailed, fmt.Errorf("failed verfication"))
	}

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		return []byte{}, esm.handleChallenge(w, r, body)
	}

	messageID := r.Header.Get("Twitch-Eventsub-Message-Id")
	if messageID == "" {
		return []byte{}, api.NewApiError(http.StatusBadRequest, fmt.Errorf("no message id"))
	}
	if _, err := esm.db.GetEventSubMessage(messageID); err == nil {
		log.Infof("Message handled before %s", messageID)
		api.WriteText(w, "handled before", http.StatusOK)
		return []byte{}, nil
	} else {
		log.Infof("Message new, handling %s", messageID)
		esm.db.CreateEventSubMessage(store.EventSubMessage{ID: messageID})
	}

	var eventSubNotification eventSubNotification

	err = json.Unmarshal(body, &eventSubNotification)
	if err != nil {
		return []byte{}, api.NewApiError(http.StatusPreconditionFailed, fmt.Errorf("failed decoding body"+err.Error()))
	}

	if eventSubNotification.Subscription.Version != "1" && eventSubNotification.Subscription.Version != "" {
		log.Errorf("Unknown subscription version found %s %s", eventSubNotification.Subscription.Version, eventSubNotification.Subscription.ID)
		return []byte{}, api.NewApiError(http.StatusOK, fmt.Errorf("unknown subscription version"))
	}

	if !esm.db.HasEventSubSubscription(eventSubNotification.Subscription.ID) {
		log.Errorf("Unknown subscription id found %s", eventSubNotification.Subscription.ID)
		return []byte{}, api.NewApiError(http.StatusOK, fmt.Errorf("unknown subscription"))
	}

	api.WriteText(w, "ok", http.StatusOK)

	return eventSubNotification.Event, nil
}

func (esm *EventsubManager) handleChallenge(w http.ResponseWriter, r *http.Request, body []byte) api.Error {
	var event struct {
		Challenge string `json:"challenge"`
	}
	err := json.Unmarshal(body, &event)
	if err != nil {
		return api.NewApiError(http.StatusBadRequest, fmt.Errorf("Failed to handle challenge: "+err.Error()))
	}

	log.Infof("Challenge success: %s", event.Challenge)
	api.WriteText(w, event.Challenge, http.StatusOK)
	return nil
}

func (esm *EventsubManager) HandleChannelPointsCustomRewardRedemption(event []byte) {
	var redemption helix.EventSubChannelPointsCustomRewardRedemptionEvent
	err := json.Unmarshal(event, &redemption)
	if err != nil {
		log.Errorf("Failed to decode event: %s", err)
		return
	}

	reward, err := esm.db.GetEnabledChannelPointRewardByID(redemption.Reward.ID)
	if err != nil {
		log.Errorf("no redemption found for rewardId %s", redemption.Reward.ID)
		return
	}

	if helixclient.RewardStatusIsUnfullfilled(redemption.Status) {
		if reward.ApproveOnly {
			if reward.Type == dto.REWARD_BTTV {
				if !esm.emoteChief.VerifyBttvRedemption(reward, redemption) {
					log.Infof("[%s] Bttv Reward did not verify refunding %s", redemption.BroadcasterUserID, redemption.Status)
					err := esm.ttlCache.Set(redemption.ID, false)
					if err != nil {
						log.Error(err)
					}
					err = esm.helixClient.UpdateRedemptionStatus(redemption.BroadcasterUserID, reward.RewardID, redemption.ID, false)
					if err != nil {
						log.Error(err)
					}
				} else {
					log.Infof("[%s] Bttv Reward is approve only, skipping redemption %s", redemption.BroadcasterUserID, redemption.Status)
					return
				}
			}
			if reward.Type == dto.REWARD_SEVENTV {
				if !esm.emoteChief.VerifySeventvRedemption(reward, redemption) {
					log.Infof("[%s] 7tv Reward did not verify refunding %s", redemption.BroadcasterUserID, redemption.Status)
					err := esm.ttlCache.Set(redemption.ID, false)
					if err != nil {
						log.Error(err)
					}
					err = esm.helixClient.UpdateRedemptionStatus(redemption.BroadcasterUserID, reward.RewardID, redemption.ID, false)
					if err != nil {
						log.Error(err)
					}
				} else {
					log.Infof("[%s] 7tv Reward is approve only, skipping redemption %s", redemption.BroadcasterUserID, redemption.Status)
					return
				}
			}
		} else {
			if reward.Type == dto.REWARD_BTTV {
				esm.emoteChief.HandleBttvRedemption(reward, redemption, true)
				return
			}
			if reward.Type == dto.REWARD_SEVENTV {
				esm.emoteChief.HandleSeventvRedemption(reward, redemption, true)
				return
			}
		}
	}
	if helixclient.RewardStatusIsCancelled(redemption.Status) {
		if reward.ApproveOnly {
			emoteID := ""
			if reward.Type == dto.REWARD_BTTV {
				emoteID, err = emotechief.GetBttvEmoteId(redemption.UserInput)
				if err != nil {
					log.Error(err)
				}
			}
			if reward.Type == dto.REWARD_SEVENTV {
				emoteID, err = emotechief.GetSevenTvEmoteId(redemption.UserInput)
				if err != nil {
					log.Error(err)
				}
			}

			if emoteID != "" {
				err := esm.db.BlockEmotes(redemption.BroadcasterUserID, []string{emoteID}, string(reward.Type))
				if err != nil {
					log.Error(err)
				}
			}
			// if we don't find the redemption in our cache, we didn't send the redemption update ourselves and need to send a rejection message
			if _, err := esm.ttlCache.Get(redemption.ID); err == ttlcache.ErrNotFound {
				esm.chatClient.Say(redemption.BroadcasterUserLogin, fmt.Sprintf("⚠️ Emote redemption by @%s was rejected", redemption.UserLogin))
			}
		}
		return
	}
	if helixclient.RewardStatusIsFullfilled(redemption.Status) {
		if reward.ApproveOnly {
			if reward.Type == dto.REWARD_BTTV {
				esm.emoteChief.HandleBttvRedemption(reward, redemption, false)
				return
			}
			if reward.Type == dto.REWARD_SEVENTV {
				esm.emoteChief.HandleSeventvRedemption(reward, redemption, false)
				return
			}
		}
	}
}

func (esm *EventsubManager) SubscribeChannelPoints(userID string) {
	response, err := esm.helixClient.CreateEventSubSubscription(userID, esm.cfg.WebhookApiBaseUrl+"/api/eventsub?type="+helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd, helix.EventSubTypeChannelPointsCustomRewardRedemptionAdd)
	if err != nil {
		log.Errorf("Error subscribing: %s", err)
		return
	}

	if response.StatusCode == http.StatusForbidden {
		log.Errorf("Forbidden subscription %s", response.ErrorMessage)
		return
	}

	log.Infof("[%d] subscription %s %s", response.StatusCode, response.Error, response.ErrorMessage)
	for _, sub := range response.Data.EventSubSubscriptions {
		log.Infof("new subscription for %s id: %s", userID, sub.ID)
		esm.db.AddEventSubSubscription(userID, sub.ID, sub.Version, sub.Type, "")
	}
}

func (esm *EventsubManager) RemoveEventSubSubscription(subscriptionID string) error {
	response, err := esm.helixClient.Client.RemoveEventSubSubscription(subscriptionID)
	if err != nil {
		return err
	}

	log.Infof("[%d] removed EventSubSubscription", response.StatusCode)
	esm.db.RemoveEventSubSubscription(subscriptionID)

	return nil
}

func (esm *EventsubManager) RemoveAllEventSubSubscriptions(userID string) {
	// @TODO rework using the DB so we don't need to query literally every sub
	resp, err := esm.helixClient.Client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{})
	if err != nil {
		log.Errorf("Failed to get subscriptions: %s", err)
		return
	}

	subscriptions := resp.Data.EventSubSubscriptions

	cursor := resp.Data.Pagination.Cursor

	for {
		if cursor == "" {
			break
		}
		log.Infof("Getting next subscriptions cursor: %s", cursor)

		nextResp, err := esm.helixClient.Client.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{After: cursor})
		if err != nil {
			log.Errorf("Failed to get subscriptions: %s", err)
		}
		cursor = nextResp.Data.Pagination.Cursor

		subscriptions = append(subscriptions, nextResp.Data.EventSubSubscriptions...)
	}

	for _, sub := range subscriptions {
		if sub.Condition.BroadcasterUserID != userID && userID != "" {
			continue
		}

		err := esm.RemoveEventSubSubscription(sub.ID)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
