package server

import (
	"net/http"

	"github.com/gempir/spamchamp/services/api/emotechief"

	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/pkg/helix"
	"github.com/gempir/spamchamp/pkg/store"
	"github.com/rs/cors"

	log "github.com/sirupsen/logrus"
)

// Server api server
type Server struct {
	cfg             *config.Config
	helixClient     *helix.Client
	helixUserClient *helix.Client
	store           *store.Store
	emotechief      *emotechief.EmoteChief
}

// NewServer create api Server
func NewServer(cfg *config.Config, helixClient *helix.Client, helixUserClient *helix.Client, store *store.Store, emotechief *emotechief.EmoteChief) Server {
	return Server{
		cfg:             cfg,
		helixClient:     helixClient,
		helixUserClient: helixUserClient,
		store:           store,
		emotechief:      emotechief,
	}
}

func (s *Server) Start() {
	go s.syncSubscriptions()
	go s.tokenRefreshRoutine()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/callback", s.handleCallback)
	mux.HandleFunc("/api/redemption", s.handleChannelPointsRedemption)
	mux.HandleFunc("/api/userConfig", s.handleUserConfig)

	handler := cors.AllowAll().Handler(mux)
	log.Info("listening on port :8035")
	err := http.ListenAndServe(":8035", handler)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}
}
