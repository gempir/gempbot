package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/humanize"
	"github.com/gempir/gempbot/internal/log"
	"github.com/nicklaw5/helix/v2"
)

type FossabotContext struct {
	Channel struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Avatar          string    `json:"avatar"`
		Slug            string    `json:"slug"`
		BroadcasterType string    `json:"broadcaster_type"`
		Provider        string    `json:"provider"`
		ProviderID      string    `json:"provider_id"`
		CreatedAt       time.Time `json:"created_at"`
		StreamTimestamp time.Time `json:"stream_timestamp"`
		IsLive          bool      `json:"is_live"`
	} `json:"channel"`
	Message struct {
		ID       string `json:"id"`
		Content  string `json:"content"`
		Provider string `json:"provider"`
		User     struct {
			ProviderID  string `json:"provider_id"`
			Login       string `json:"login"`
			DisplayName string `json:"display_name"`
			Roles       []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Type string `json:"type"`
			} `json:"roles"`
		} `json:"user"`
	} `json:"message"`
}

const PREDICTION_LOCK = "lock"
const PREDICTION_CANCEL = "cancel"

func (a *Api) FossabotHandler(w http.ResponseWriter, r *http.Request) {
	customApiToken := r.Header.Get("x-fossabot-customapitoken")

	resp, err := http.Get("https://api.fossabot.com/v2/customapi/context/" + customApiToken)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	var context FossabotContext
	if err := json.NewDecoder(resp.Body).Decode(&context); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if strings.HasPrefix(context.Message.Content, "!prediction") {
		a.HandlePrediction(context, w, r)
		return
	}

	http.Error(w, "no relevant command", http.StatusNoContent)
}

func (a *Api) HandlePrediction(context FossabotContext, w http.ResponseWriter, r *http.Request) {
	var responseMsg string

	if strings.HasPrefix(context.Message.Content, "!prediction start") {
		responseMsg = a.startPrediction(context)
	} else if strings.HasPrefix(context.Message.Content, "!prediction end") {
		responseMsg = a.setOutcomeForPrediction(context)
	} else if strings.HasPrefix(context.Message.Content, "!prediction cancel") {
		responseMsg = a.lockOrCancelPrediction(context, PREDICTION_CANCEL)
	} else if strings.HasPrefix(context.Message.Content, "!prediction lock") {
		responseMsg = a.lockOrCancelPrediction(context, PREDICTION_LOCK)
	}

	api.WriteText(w, responseMsg, http.StatusOK)
}

func (h *Api) lockOrCancelPrediction(context FossabotContext, status string) string {
	resp, err := h.helixClient.GetPredictions(&helix.PredictionsParams{BroadcasterID: context.Channel.ProviderID})
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	prediction := resp.Data.Predictions[0]

	token, err := h.db.GetUserAccessToken(context.Channel.ProviderID)
	if err != nil {
		return err.Error()
	}
	h.helixClient.SetUserAccessToken(token.AccessToken)
	resp, err = h.helixClient.EndPrediction(&helix.EndPredictionParams{BroadcasterID: context.Channel.ProviderID, ID: prediction.ID, Status: status})
	h.helixClient.SetUserAccessToken("")

	if err != nil {
		log.Error(err)
		return err.Error()
	}
	log.Infof("[helix] %d CancelOrLockPrediction %s", resp.StatusCode, context.Channel.ProviderID)
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Sprintf("Bad Twitch API response %d", resp.StatusCode)
	}

	return ""
}

func (h *Api) setOutcomeForPrediction(context FossabotContext) string {
	input := strings.TrimPrefix(context.Message.Content, "!prediction end ")

	var winningOutcome helix.Outcomes

	resp, err := h.helixClient.GetPredictions(&helix.PredictionsParams{BroadcasterID: context.Channel.ProviderID})
	if err != nil {
		log.Error(err)
		return err.Error()
	}
	prediction := resp.Data.Predictions[0]

	for index, outcome := range prediction.Outcomes {
		if strings.EqualFold(outcome.Title, input) || fmt.Sprintf("%d", index+1) == input {
			winningOutcome = outcome
			break
		}
	}

	if winningOutcome.ID == "" {
		return "outcome not found"
	}

	_, err = h.helixClient.EndPrediction(&helix.EndPredictionParams{BroadcasterID: context.Channel.ProviderID, ID: prediction.ID, Status: dto.PredictionStatusResolved, WinningOutcomeID: winningOutcome.ID})
	if err != nil {
		log.Error(err)
		return fmt.Sprintf("failed to end prediction: %s", err.Error())
	}

	return ""
}

func (h *Api) startPrediction(context FossabotContext) string {
	input := strings.TrimPrefix(context.Message.Content, "!prediction start ")

	split := strings.Split(input, ";")

	if len(split) < 1 {
		return "missing title"
	}

	title := strings.TrimSpace(split[0])
	predictionWindow := 60

	if len(split) >= 2 {
		var err error
		predictionWindow, err = humanize.StringToSeconds(strings.TrimSpace(split[1]))
		if err != nil {
			log.Error(err)
			return "failed to parse time"
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
		BroadcasterID:    context.Channel.ProviderID,
		Title:            title,
		Outcomes:         outcomes,
		PredictionWindow: predictionWindow,
	}

	_, err := h.helixClient.CreatePrediction(prediction)
	if err != nil {
		log.Error(err)
		return fmt.Sprintf("failed to create prediction: %s", err.Error())
	}

	return ""
}
