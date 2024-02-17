package ysweet

import (
	"context"
	"strings"

	"github.com/carlmjohnson/requests"
	"github.com/gempir/gempbot/internal/config"
)

type Factory struct {
	ysweetUrl   string
	ssl         bool
	bearerToken string
}

func NewFactory(cfg *config.Config) *Factory {
	return &Factory{
		ysweetUrl:   cfg.YsweetUrl,
		ssl:         strings.HasPrefix(cfg.WebBaseUrl, "https"),
		bearerToken: cfg.YsweetToken,
	}
}

type TokenResponse struct {
	Url   string `json:"url"`
	DocId string `json:"docId"`
	Token string `json:"token"`
}

type DocResponse struct {
	DocID string `json:"docId"`
}

func (f *Factory) CreateToken(docID string) (TokenResponse, error) {
	var docResponse DocResponse
	err := requests.
		URL(f.ysweetUrl).
		Post().
		BodyJSON(map[string]string{"docId": docID}).
		Bearer(f.bearerToken).
		Pathf("/doc/new").
		ToJSON(&docResponse).
		Fetch(context.Background())
	if err != nil {
		return TokenResponse{}, err
	}

	var tokenResponse TokenResponse
	err = requests.
		URL(f.ysweetUrl).
		Post().
		Pathf("/doc/%s/auth", docResponse.DocID).
		BodyJSON(map[string]string{}).
		Bearer(f.bearerToken).
		ToJSON(&tokenResponse).
		Fetch(context.Background())
	if err != nil {
		return TokenResponse{}, err
	}

	if f.ssl {
		tokenResponse.Url = strings.Replace(tokenResponse.Url, "ws://", "wss://", 1)
	}

	return tokenResponse, nil
}
