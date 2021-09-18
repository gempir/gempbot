package emotechief

import (
	"net/http"

	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/store"
)

type EmoteChief struct {
	store      *store.Redis
	cfg        *config.Config
	db         *store.Database
	httpClient *http.Client
}

func NewEmoteChief(store *store.Redis, cfg *config.Config, db *store.Database) *EmoteChief {
	return &EmoteChief{
		store:      store,
		cfg:        cfg,
		db:         db,
		httpClient: &http.Client{},
	}
}
