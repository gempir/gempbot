package emotechief

import (
	"net/http"

	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/pkg/store"
)

type EmoteChief struct {
	store      *store.Store
	cfg        *config.Config
	httpClient *http.Client
}

func NewEmoteChief(store *store.Store, cfg *config.Config) *EmoteChief {
	return &EmoteChief{
		store:      store,
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}
