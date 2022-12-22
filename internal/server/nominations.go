package server

import (
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
)

func (a *Api) NominationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	pageSize := 10
	page := 1
	if r.URL.Query().Get("page") != "" {
		page = 1
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

	election, err := a.db.GetActiveElection(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "no active election", http.StatusBadRequest)
		return
	}

	nominations, err := a.db.GetNominations(r.Context(), user.ID, election.ID, page, pageSize)
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
