package server

import (
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

func (a *Api) NominationsBlockHandler(w http.ResponseWriter, r *http.Request) {
	authResp, _, apiErr := a.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID

	if r.Method == http.MethodDelete {
		if r.URL.Query().Get("managing") != "" {
			userID, apiErr = a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
			if apiErr != nil {
				http.Error(w, apiErr.Error(), apiErr.Status())
				return
			}
		}

		emoteID := r.URL.Query().Get("emoteID")

		err := a.db.ClearNominationEmote(r.Context(), userID, emoteID)
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

func (a *Api) NominationVoteHandler(w http.ResponseWriter, r *http.Request) {
	authResp, _, apiErr := a.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID

	channel, err := a.helixClient.GetUserByUsername(r.URL.Query().Get("channel"))
	if err != nil {
		http.Error(w, "user not found", http.StatusBadRequest)
		return
	}

	emoteID := r.URL.Query().Get("emoteID")

	if r.Method == http.MethodPost {
		count, err := a.db.CountNominationVotes(r.Context(), channel.ID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		election, err := a.db.GetElection(r.Context(), channel.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count >= election.VoteAmount {
			api.WriteJson(w, "max votes reached", http.StatusBadRequest)
			return
		}

		err = a.db.CreateNominationVote(r.Context(), store.NominationVote{EmoteID: emoteID, ChannelTwitchID: channel.ID, VoteBy: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		api.WriteJson(w, "ok", http.StatusOK)
	}
	if r.Method == http.MethodDelete {
		nom, err := a.db.GetNomination(r.Context(), channel.ID, emoteID)
		if err == nil {
			if nom.NominatedBy == userID {
				err = a.db.RemoveNomination(r.Context(), userID, emoteID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				api.WriteJson(w, "ok", http.StatusOK)
				return
			}
		}

		err = a.db.RemoveNominationVote(r.Context(), store.NominationVote{EmoteID: emoteID, ChannelTwitchID: channel.ID, VoteBy: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		api.WriteJson(w, "ok", http.StatusOK)
	}

	http.Error(w, "", http.StatusMethodNotAllowed)
}

func (a *Api) NominationDownvoteHandler(w http.ResponseWriter, r *http.Request) {
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

	emoteID := r.URL.Query().Get("emoteID")

	if r.Method == http.MethodPost {
		count, err := a.db.CountNominationDownvotes(r.Context(), user.ID, userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		election, err := a.db.GetElection(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count >= election.VoteAmount {
			api.WriteJson(w, "max downvotes reached", http.StatusBadRequest)
			return
		}

		nom, err := a.db.GetNomination(r.Context(), user.ID, emoteID)
		if err == nil {
			if nom.NominatedBy == userID {
				err = a.db.RemoveNomination(r.Context(), userID, emoteID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				api.WriteJson(w, "ok", http.StatusOK)
				return
			}
		}

		err = a.db.CreateNominationDownvote(r.Context(), store.NominationDownvote{EmoteID: emoteID, ChannelTwitchID: user.ID, VoteBy: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		api.WriteJson(w, "ok", http.StatusOK)
	}
	if r.Method == http.MethodDelete {
		err = a.db.RemoveNominationDownvote(r.Context(), store.NominationDownvote{EmoteID: emoteID, ChannelTwitchID: user.ID, VoteBy: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		api.WriteJson(w, "ok", http.StatusOK)
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
