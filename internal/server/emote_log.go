package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/log"
)

func (a *Api) EmoteLogHandler(w http.ResponseWriter, r *http.Request) {

	limit := 5

	if r.URL.Query().Get("limit") != "" {
		limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	}

	if r.URL.Query().Get("channel") == "" {
		api.WriteJson(w, fmt.Errorf("channel is required"), http.StatusBadRequest)
		return
	}

	channel, err := a.helixClient.GetUserByUsername(r.URL.Query().Get("channel"))
	if err != nil {
		api.WriteJson(w, fmt.Errorf("channel not found"), http.StatusBadRequest)
		return
	}

	entries := a.db.GetEmoteLogEntries(r.Context(), channel.ID, limit)

	users := []string{}
	for _, entry := range entries {
		users = append(users, entry.AddedBy)
	}

	usersMap, err := a.helixClient.GetUsersByUserIds(users)
	if err != nil {
		log.Error(err.Error())
	}

	for i, entry := range entries {
		if user, ok := usersMap[entry.AddedBy]; ok {
			entries[i].AddedBy = user.DisplayName
		}
	}

	if r.URL.Query().Has("text") {
		logs := []string{}
		for _, entry := range entries {
			typeAdd := "🗳️"
			if entry.Type == dto.REWARD_SEVENTV {
				typeAdd = "🔄"
			}

			logs = append(logs, fmt.Sprintf("%s %s [] by %s", typeAdd, entry.EmoteCode, entry.AddedBy))
		}
		api.WriteText(w, strings.Join(logs, ", "), http.StatusOK)

		return
	}

	api.WriteJson(w, entries, http.StatusOK)
}
