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
	prefixStripped := strings.TrimPrefix(payload.Msg.Message, "!prediction")
	split := strings.Split(prefixStripped, ";")

	if len(split) <= 1 {
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
		h.handleError(payload.Msg, fmt.Errorf("bad twitch api response %d", resp.StatusCode))
		return
	}

	h.handleSuccess(payload.Msg, *prediction)
}

func (h *Handler) handleError(msg twitch.PrivateMessage, err error) {
	h.redis.PublishSpeakerMessage(msg.Channel, fmt.Sprintf("@%s failed to create prediction: %s", msg.User.DisplayName, err))
}

func (h *Handler) handleSuccess(msg twitch.PrivateMessage, prediction nickHelix.CreatePredictionParams) {
	h.redis.PublishSpeakerMessage(
		msg.Channel,
		fmt.Sprintf("PogChamp New prediction: %s - %s | %s - ending in %s created by @%s",
			prediction.Title,
			prediction.Outcomes[0].Title,
			prediction.Outcomes[1].Title,
			humanize.SecondsToString(prediction.PredictionWindow),
			msg.User.DisplayName,
		),
	)
}
