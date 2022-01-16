package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/log"
)

func (a *Api) EmoteHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID := ""
	login := ""

	authResult, _, err := a.authClient.AttemptAuth(r, w)
	if err != nil {
		return
	}
	userID = authResult.Data.UserID
	login = authResult.Data.Login

	if r.URL.Query().Get("managing") != "" {
		userID, err := a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
		if err != nil {
			http.Error(w, err.Error(), err.Status())
			return
		}

		uData, helixError := a.helixClient.GetUserByUserID(userID)
		if helixError != nil {
			api.WriteJson(w, fmt.Errorf("could not find managing user in helix"), http.StatusBadRequest)
			return
		}
		login = uData.Login
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

			emote, err := a.emoteChief.RemoveSevenTvEmote(userID, emoteID)
			if err != nil || emote == nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			a.bot.ChatClient.Say(login, fmt.Sprintf("⚠️ Emote %s has been removed and blocked", emote.Name))
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

			a.bot.ChatClient.Say(login, fmt.Sprintf("⚠️ Emote %s has been removed and blocked", emote.Code))
		}

		api.WriteJson(w, "ok", http.StatusOK)
		return
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageNumber, convError := strconv.Atoi(page)
	if convError != nil {
		http.Error(w, convError.Error(), http.StatusBadRequest)
		return
	}

	api.WriteJson(w, a.db.GetEmoteHistory(r.Context(), userID, pageNumber, api.EMOTEHISTORY_PAGE_SIZE, r.URL.Query().Has("added")), http.StatusOK)
}
