package predictions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gempir/bitraft/pkg/dto"
	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/humanize"
	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/gempir/go-twitch-irc/v2"
	nickHelix "github.com/nicklaw5/helix"
)

type Handler struct {
	redis       *store.Redis
	db          *store.Database
	helixClient *helix.Client
}

func NewHandler(helixClient *helix.Client, redis *store.Redis, db *store.Database) *Handler {
	return &Handler{
		redis:       redis,
		db:          db,
		helixClient: helixClient,
	}
}

// !prediction Will nymn win this game?;yes;no;3m --> yes;no;3m
// !prediction Will he win                        --> yes;no;1m
// !prediction Will he win;maybe                  --> maybe;no;1m
func (h *Handler) HandleCommand(payload dto.CommandPayload) {
	switch payload.Name {
	case dto.CmdNameOutcome:
		h.setOutcomeForPrediction(payload)
	case dto.CmdNamePrediction:
		h.startPrediction(payload)
	}
}

func (h *Handler) setOutcomeForPrediction(payload dto.CommandPayload) {
	var winningOutcome store.PredictionLogOutcome

	prediction, err := h.db.GetActivePrediction(payload.Msg.RoomID)
	if err != nil {
		h.handleError(payload.Msg, errors.New("no active prediction found"))
		return
	}

	outcomes := h.db.GetOutcomes(prediction.ID)
	for _, outcome := range outcomes {
		if payload.Query == "1" || payload.Query == dto.Outcome_First || payload.Query == "first" {
			if outcome.Color == dto.Outcome_First {
				winningOutcome = outcome
			}
		}
		if payload.Query == "2" || payload.Query == dto.Outcome_Second || payload.Query == "red" || payload.Query == "second" {
			if outcome.Color == dto.Outcome_Second {
				winningOutcome = outcome
			}
		}
	}

	if winningOutcome.ID == "" {
		h.handleError(payload.Msg, errors.New("outcome not found"))
		return
	}

	token, err := h.db.GetUserAccessToken(payload.Msg.RoomID)
	if err != nil {
		h.handleError(payload.Msg, errors.New("no api token, broadcaster needs to login again in dashboard"))
		return
	}

	h.helixClient.Client.SetUserAccessToken(token.AccessToken)
	resp, err := h.helixClient.Client.EndPrediction(&nickHelix.EndPredictionParams{BroadcasterID: payload.Msg.RoomID, ID: prediction.ID, Status: dto.PredictionStatusResolved, WinningOutcomeID: winningOutcome.ID})
	h.helixClient.Client.SetUserAccessToken("")

	if err != nil {
		log.Error(err)
		h.handleError(payload.Msg, errors.New("bad twitch api response"))
		return
	}
	log.Infof("[helix] %d EndPrediction %s", resp.StatusCode, payload.Msg.RoomID)
	if resp.StatusCode >= http.StatusBadRequest {
		h.handleError(payload.Msg, fmt.Errorf("bad twitch api response %s", resp.ErrorMessage))
		return
	}

	for _, prediction := range resp.Data.Predictions {
		err := h.db.SavePrediction(store.PredictionLog{ID: prediction.ID, OwnerTwitchID: payload.Msg.RoomID, Title: prediction.Title, StartedAt: prediction.CreatedAt.Time, LockedAt: prediction.LockedAt.Time, EndedAt: prediction.EndedAt.Time, WinningOutcomeID: prediction.WinningOutcomeID})
		if err != nil {
			log.Error(err)
		}

		for _, outcome := range prediction.Outcomes {
			err := h.db.SaveOutcome(store.PredictionLogOutcome{ID: outcome.ID, PredictionID: prediction.ID, Title: outcome.Title, Color: outcome.Color, Users: outcome.Users, ChannelPoints: outcome.ChannelPoints})
			if err != nil {
				log.Error(err)
			}
		}
	}

	h.redis.PublishSpeakerMessage(
		payload.Msg.Channel,
		fmt.Sprintf("PogChamp ended prediction: %s Winner: %s %s",
			prediction.Title,
			winningOutcome.GetColorEmoji(),
			winningOutcome.Title,
		),
	)
}

func (h *Handler) startPrediction(payload dto.CommandPayload) {
	split := strings.Split(payload.Query, ";")

	if len(split) < 1 {
		h.handleError(payload.Msg, errors.New("no title given"))
		return
	}

	title := strings.TrimSpace(split[0])
	outcome1 := "yes"
	outcome2 := "no"
	predictionWindow := 60

	if len(split) >= 2 {
		outcome1 = strings.TrimSpace(split[1])
	}
	if len(split) >= 3 {
		outcome2 = strings.TrimSpace(split[2])
	}
	if len(split) >= 4 {
		var err error
		predictionWindow, err = humanize.StringToSeconds(strings.TrimSpace(split[3]))
		if err != nil {
			log.Error(err)
			h.handleError(payload.Msg, errors.New("failed to parse time"))
			return
		}
	}

	if predictionWindow > 1800 {
		h.handleError(payload.Msg, errors.New("max 30 minutes"))
		return
	}

	prediction := &nickHelix.CreatePredictionParams{
		BroadcasterID:    payload.Msg.RoomID,
		Title:            title,
		Outcomes:         []nickHelix.PredictionChoiceParam{{Title: outcome1}, {Title: outcome2}},
		PredictionWindow: predictionWindow,
	}

	token, err := h.db.GetUserAccessToken(payload.Msg.RoomID)
	if err != nil {
		h.handleError(payload.Msg, errors.New("no api token, broadcaster needs to login again in dashboard"))
		return
	}

	h.helixClient.Client.SetUserAccessToken(token.AccessToken)
	resp, err := h.helixClient.Client.CreatePrediction(prediction)
	h.helixClient.Client.SetUserAccessToken("")

	if err != nil {
		log.Error(err)
		h.handleError(payload.Msg, errors.New("bad twitch api response"))
		return
	}
	log.Infof("[helix] %d CreatePrediction %s", resp.StatusCode, payload.Msg.RoomID)
	if resp.StatusCode >= http.StatusBadRequest {
		h.handleError(payload.Msg, fmt.Errorf("bad twitch api response: %s", resp.ErrorMessage))
		return
	}

	for _, prediction := range resp.Data.Predictions {
		err := h.db.SavePrediction(store.PredictionLog{ID: prediction.ID, OwnerTwitchID: payload.Msg.RoomID, Title: prediction.Title, StartedAt: prediction.CreatedAt.Time, LockedAt: prediction.LockedAt.Time, EndedAt: prediction.EndedAt.Time, WinningOutcomeID: prediction.WinningOutcomeID})
		if err != nil {
			log.Error(err)
		}

		for _, outcome := range prediction.Outcomes {
			err := h.db.SaveOutcome(store.PredictionLogOutcome{ID: outcome.ID, PredictionID: prediction.ID, Title: outcome.Title, Color: strings.ToLower(outcome.Color), Users: outcome.Users, ChannelPoints: outcome.ChannelPoints})
			if err != nil {
				log.Error(err)
			}
		}
	}

	h.redis.PublishSpeakerMessage(
		payload.Msg.Channel,
		fmt.Sprintf("PogChamp prediction: %s [ %s | %s ] ending in %s",
			prediction.Title,
			prediction.Outcomes[0].Title,
			prediction.Outcomes[1].Title,
			humanize.SecondsToString(prediction.PredictionWindow),
		),
	)
}

func (h *Handler) handleError(msg twitch.PrivateMessage, err error) {
	h.redis.PublishSpeakerMessage(msg.Channel, fmt.Sprintf("@%s %s", msg.User.DisplayName, err))
}
