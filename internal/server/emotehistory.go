package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
)

func (a *Api) EmoteHistoryHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	userID := ""

	if username == "" {
		authResult, _, err := a.authClient.AttemptAuth(r, w)
		if err != nil {
			return
		}
		userID = authResult.Data.UserID

		if r.URL.Query().Get("managing") != "" {
			userID, err = a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
			if err != nil {
				http.Error(w, err.Error(), err.Status())
				return
			}
		}
	} else {
		user, err := a.helixClient.GetUserByUsername(username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID = user.ID
	}
	if r.Method == http.MethodDelete {
		a.db.RemoveEmoteAdd(userID, r.URL.Query().Get("emoteId"))

		api.WriteJson(w, "ok", http.StatusOK)
		return
	}
	if r.Method == http.MethodPatch {
		emoteID := r.URL.Query().Get("emoteId")

		a.db.BlockEmoteAdd(userID, emoteID)
		emoteAdd := a.db.GetEmoteAdd(userID, emoteID)

		if emoteAdd.Type == dto.REWARD_SEVENTV {
			err := a.db.BlockEmotes(userID, []string{emoteID}, string(dto.REWARD_SEVENTV))
			if err != nil {
				log.Error(err)
			}

			err = a.sevenTvClient.RemoveEmote(userID, emoteID)
			a.db.CreateEmoteAdd(userID, dto.REWARD_SEVENTV, emoteID, dto.EMOTE_ADD_REMOVED_BLOCKED)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			emote, err := a.sevenTvClient.GetEmote(emoteID)
			if err != nil {
				log.Error(err)
				return
			}

			a.bot.Send(userID, fmt.Sprintf("⚠️ Emote %s has been removed and blocked", emote.Code))
		} else if emoteAdd.Type == dto.REWARD_BTTV {
			err := a.db.BlockEmotes(userID, []string{emoteID}, string(dto.REWARD_BTTV))
			if err != nil {
				log.Error(err)
			}

			emote, err := a.emoteChief.RemoveBttvEmote(userID, emoteID)
			if err != nil || emote == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			a.bot.Send(userID, fmt.Sprintf("⚠️ Emote %s has been removed and blocked", emote.Code))
		}

		api.WriteJson(w, "ok", http.StatusOK)
		return
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	api.WriteJson(w, a.db.GetEmoteHistory(r.Context(), userID, pageNumber, 20, r.URL.Query().Has("added")), http.StatusOK)
}
