package emotechief

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/store"
)

type EmoteChief struct {
	cfg        *config.Config
	db         *store.Database
	httpClient *http.Client
}

func NewEmoteChief(cfg *config.Config, db *store.Database) *EmoteChief {
	return &EmoteChief{
		cfg:        cfg,
		db:         db,
		httpClient: &http.Client{},
	}
}
