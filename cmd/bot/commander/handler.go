package commander

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gempir/gempbot/cmd/bot/scaler"
	"github.com/gempir/gempbot/pkg/apiclient"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/humanize"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/go-twitch-irc/v2"
	nickHelix "github.com/nicklaw5/helix/v2"
)

type Handler struct {
	cfg         *config.Config
	apiClient   *apiclient.ApiClient
	db          *store.Database
	helixClient *helix.Client
	write       chan scaler.Message
}

func NewHandler(cfg *config.Config, helixClient *helix.Client, db *store.Database, write chan scaler.Message) *Handler {
	return &Handler{
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
		apiClient:   apiclient.NewApiClient(cfg),
		write:       write,
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
		h.handlePrediction(payload)
	}
}

func (h *Handler) handlePrediction(payload dto.CommandPayload) {
	if strings.ToLower(payload.Query) == "lock" {
		h.lockOrCancelPrediction(payload, dto.PredictionStatusLocked)
		return
	}
	if strings.ToLower(payload.Query) == "cancel" {
		h.lockOrCancelPrediction(payload, dto.PredictionStatusCanceled)
		return
	}

	h.startPrediction(payload)
}

func (h *Handler) lockOrCancelPrediction(payload dto.CommandPayload, status string) {
	prediction, err := h.db.GetActivePrediction(payload.Msg.RoomID)
	if err != nil {
		h.handleError(payload.Msg, errors.New("no active prediction found"))
		return
	}

	token, err := h.db.GetUserAccessToken(payload.Msg.RoomID)
	if err != nil {
		h.handleError(payload.Msg, errors.New("no api token, broadcaster needs to login again in dashboard"))
		return
	}
	h.helixClient.Client.SetUserAccessToken(token.AccessToken)
	resp, err := h.helixClient.Client.EndPrediction(&nickHelix.EndPredictionParams{BroadcasterID: payload.Msg.RoomID, ID: prediction.ID, Status: status})
	h.helixClient.Client.SetUserAccessToken("")

	if err != nil {
		log.Error(err)
		h.handleError(payload.Msg, errors.New("bad twitch api response"))
		return
	}
	log.Infof("[helix] %d CancelOrLockPrediction %s", resp.StatusCode, payload.Msg.RoomID)
	if resp.StatusCode >= http.StatusBadRequest {
		h.handleError(payload.Msg, fmt.Errorf("bad twitch api response %s", resp.ErrorMessage))
		return
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
		err := h.db.SavePrediction(store.PredictionLog{ID: prediction.ID, OwnerTwitchID: payload.Msg.RoomID, Title: prediction.Title, StartedAt: prediction.CreatedAt.Time, LockedAt: &prediction.LockedAt.Time, EndedAt: &prediction.EndedAt.Time, WinningOutcomeID: prediction.WinningOutcomeID})
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

	resp, err := h.apiClient.CreatePrediction(payload.Msg.RoomID, prediction)
	if err != nil {
		log.Error(err)
		h.handleError(payload.Msg, errors.New(err.Error()))
		return
	}
	log.Infof("[api] %d CreatePrediction %s", resp.StatusCode, payload.Msg.RoomID)
}

func (h *Handler) handleError(msg twitch.PrivateMessage, err error) {
	h.write <- scaler.Message{Channel: msg.Channel, Message: fmt.Sprintf("@%s %s", msg.User.DisplayName, err)}
}
