package emotechief

import (
	"net/http"

	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/store"
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
