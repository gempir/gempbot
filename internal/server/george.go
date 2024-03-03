package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type GeorgeRequest struct {
	Query    string `json:"query"`
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Month    int    `json:"month"`
	Year     int    `json:"year"`
	Day      int    `json:"day"`
	Model    string `json:"model"`
	Limit    int    `json:"limit"`
}

var allowlistedUsers = map[string]string{
	"gempir":      "gempir",
	"fawcan":      "fawcan",
	"xanabilek":   "xanabilek",
	"mr0lle":      "mr0lle",
	"hotbear1110": "hotbear1110",
	"leppunen":    "leppunen",
	"nymn":        "nymn",
	"yabbe":       "yabbe",
	"pajlada":     "pajlada",
}

var mutex = sync.Mutex{}

func (a *Api) GeorgeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	authResp, _, apiErr := a.authClient.AttemptAuth(r, w)
	if apiErr != nil {
		return
	}
	if _, ok := allowlistedUsers[authResp.Data.Login]; !ok {
		http.Error(w, "Not allowedlisted", http.StatusForbidden)
		return
	}

	if mutex.TryLock() {
		defer mutex.Unlock()
	} else {
		http.Error(w, "Already processing", http.StatusTooManyRequests)
		return
	}

	var req GeorgeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	err = a.george.AnalyzeUser(req.Query, req.Channel, req.Username, req.Month, req.Year, req.Day, req.Model, req.Limit, r.Context(), func(chunk string) {
		_, err := fmt.Fprint(w, chunk)
		if err != nil {
			fmt.Println("Error writing response:", err)
			return
		}
		w.(http.Flusher).Flush()
	})

	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			http.Error(w, "gempir is offline, his GPU needs to be connected", http.StatusInternalServerError)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
