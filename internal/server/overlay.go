package server

import (
	"net/http"
	"strings"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/store"
	"github.com/gempir/gempbot/internal/ysweet"
	"github.com/google/uuid"
	"github.com/teris-io/shortid"
)

type OverlayResponse struct {
	Overlay store.Overlay        `json:"overlay"`
	Auth    ysweet.TokenResponse `json:"auth"`
}

func (a *Api) OverlayHandler(w http.ResponseWriter, r *http.Request) {
	authResp, _, apiErr := a.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	userID := authResp.Data.UserID

	if r.URL.Query().Get("managing") != "" {
		userID, apiErr = a.userAdmin.CheckEditor(r, a.userAdmin.GetUserConfig(userID))
		if apiErr != nil {
			http.Error(w, apiErr.Error(), apiErr.Status())
			return
		}
	}

	if r.Method == http.MethodGet {
		if r.URL.Query().Get("roomId") != "" {
			overlay := a.db.GetOverlayByRoomId(r.URL.Query().Get("roomId"))

			token, err := a.tokenFactory.CreateToken(overlay.RoomID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			api.WriteJson(w, OverlayResponse{overlay, token}, http.StatusOK)
			return
		}

		if r.URL.Query().Get("id") != "" {
			overlay := a.db.GetOverlay(r.URL.Query().Get("id"), userID)
			token, err := a.tokenFactory.CreateToken(overlay.RoomID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			api.WriteJson(w, OverlayResponse{overlay, token}, http.StatusOK)
			return
		}

		overlays := a.db.GetOverlays(userID)
		api.WriteJson(w, overlays, http.StatusOK)
	} else if r.Method == http.MethodPost {
		overlay := store.Overlay{}
		overlay.OwnerTwitchID = userID
		overlay.ID = shortid.MustGenerate()
		// long string so you cant read addressbar easily
		var roomID []string
		for i := 0; i < 4; i++ {
			roomID = append(roomID, uuid.New().String())
		}
		overlay.RoomID = strings.Join(roomID, "-")

		err := a.db.SaveOverlay(overlay)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		api.WriteJson(w, overlay, http.StatusCreated)

	} else if r.Method == http.MethodDelete {
		if r.URL.Query().Get("id") == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
		}

		a.db.DeleteOverlay(r.URL.Query().Get("id"))
		w.WriteHeader(http.StatusOK)
	}
}
