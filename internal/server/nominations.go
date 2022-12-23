package server

import (
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

func (a *Api) NominationVoteHandler(w http.ResponseWriter, r *http.Request) {
	authResp, _, apiErr := a.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID

	user, err := a.helixClient.GetUserByUsername(r.URL.Query().Get("channel"))
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		err = a.db.CreateNominationVote(r.Context(), store.NominationVote{EmoteID: r.URL.Query().Get("emoteID"), ChannelTwitchID: user.ID, VoteBy: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		api.WriteJson(w, "ok", http.StatusOK)
	}
	if r.Method == http.MethodDelete {
		if r.URL.Query().Get("managing") != "" {
			userID, apiErr = a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
			if apiErr != nil {
				http.Error(w, apiErr.Error(), apiErr.Status())
				return
			}
		}

		emoteID := r.URL.Query().Get("emoteID")

		err = a.db.ClearNominationEmote(r.Context(), userID, emoteID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = a.db.BlockEmotes(userID, []string{emoteID}, string(dto.REWARD_SEVENTV))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		api.WriteJson(w, "ok", http.StatusOK)
		return
	}

	http.Error(w, "", http.StatusMethodNotAllowed)
}

func (a *Api) NominationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	channel := r.URL.Query().Get("channel")
	if channel == "" {
		http.Error(w, "no channel given", http.StatusBadRequest)
		return
	}
	user, err := a.helixClient.GetUserByUsername(channel)
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	_, err = a.db.GetActiveElection(r.Context(), user.ID)
	if err != nil {
		api.WriteJson(w, []store.Nomination{}, http.StatusOK)
		return
	}

	nominations, err := a.db.GetNominations(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ids := []string{}
	for _, nomination := range nominations {
		ids = append(ids, nomination.NominatedBy)
	}

	users, err := a.helixClient.GetUsersByUserIds(ids)
	if err != nil {
		log.Errorf("Failed to fetch users %s", err.Error())
	}

	transformedNominations := []store.Nomination{}
	for _, nomination := range nominations {
		if _, ok := users[nomination.NominatedBy]; ok {
			nomination.NominatedBy = users[nomination.NominatedBy].DisplayName
		}
		transformedNominations = append(transformedNominations, nomination)
	}

	api.WriteJson(w, transformedNominations, http.StatusOK)
}
