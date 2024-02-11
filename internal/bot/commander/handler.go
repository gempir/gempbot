package commander

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/humanize"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/nicklaw5/helix/v2"
)

type Handler struct {
	cfg         *config.Config
	db          *store.Database
	helixClient helixclient.Client
	chatSay     func(channelID, message string)
}

func NewHandler(cfg *config.Config, helixClient helixclient.Client, db *store.Database, chatSay func(channelID, message string)) *Handler {
	return &Handler{
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
		chatSay:     chatSay,
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
	resp, err := h.helixClient.GetPredictions(&helix.PredictionsParams{BroadcasterID: payload.Msg.RoomID})
	if err != nil {
		log.Error(err)
		h.handleError(payload.Msg, err)
		return
	}
	prediction := resp.Data.Predictions[0]

	token, err := h.db.GetUserAccessToken(payload.Msg.RoomID)
	if err != nil {
		h.handleError(payload.Msg, errors.New("no api token, broadcaster needs to login again in dashboard"))
		return
	}
	h.helixClient.SetUserAccessToken(token.AccessToken)
	resp, err = h.helixClient.EndPrediction(&helix.EndPredictionParams{BroadcasterID: payload.Msg.RoomID, ID: prediction.ID, Status: status})
	h.helixClient.SetUserAccessToken("")

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
	var winningOutcome helix.Outcomes

	resp, err := h.helixClient.GetPredictions(&helix.PredictionsParams{BroadcasterID: payload.Msg.RoomID})
	if err != nil {
		log.Error(err)
		h.handleError(payload.Msg, err)
		return
	}
	prediction := resp.Data.Predictions[0]

	for index, outcome := range prediction.Outcomes {
		if strings.EqualFold(outcome.Title, payload.Query) || fmt.Sprintf("%d", index+1) == payload.Query {
			winningOutcome = outcome
			break
		}
	}

	if winningOutcome.ID == "" {
		h.handleError(payload.Msg, errors.New("outcome not found"))
		return
	}

	_, err = h.helixClient.EndPrediction(&helix.EndPredictionParams{BroadcasterID: payload.Msg.RoomID, ID: prediction.ID, Status: dto.PredictionStatusResolved, WinningOutcomeID: winningOutcome.ID})
	if err != nil {
		log.Error(err)
		h.handleError(payload.Msg, errors.New("bad twitch api response"))
		return
	}
}

func (h *Handler) startPrediction(payload dto.CommandPayload) {
	split := strings.Split(payload.Query, ";")

	if len(split) < 1 {
		h.handleError(payload.Msg, errors.New("no title given"))
		return
	}

	title := strings.TrimSpace(split[0])
	predictionWindow := 60

	if len(split) >= 2 {
		var err error
		predictionWindow, err = humanize.StringToSeconds(strings.TrimSpace(split[1]))
		if err != nil {
			log.Error(err)
			h.handleError(payload.Msg, errors.New("failed to parse time"))
			return
		}
	}

	outcomes := []helix.PredictionChoiceParam{}
	if len(split) >= 3 {
		for _, outcome := range split[2:] {
			outcomes = append(outcomes, helix.PredictionChoiceParam{
				Title: outcome,
			})
		}
	}
	if len(outcomes) == 0 {
		outcomes = append(outcomes, helix.PredictionChoiceParam{
			Title: "yes",
		})
	}
	if len(outcomes) == 1 {
		outcomes = append(outcomes, helix.PredictionChoiceParam{
			Title: "no",
		})
	}

	prediction := &helix.CreatePredictionParams{
		BroadcasterID:    payload.Msg.RoomID,
		Title:            title,
		Outcomes:         outcomes,
		PredictionWindow: predictionWindow,
	}

	_, err := h.helixClient.CreatePrediction(prediction)
	if err != nil {
		log.Error(err)
		h.handleError(payload.Msg, errors.New(err.Error()))
		return
	}
}

func (h *Handler) handleError(msg twitch.PrivateMessage, err error) {
	h.chatSay(msg.RoomID, fmt.Sprintf("@%s %s", msg.User.DisplayName, err))
}
