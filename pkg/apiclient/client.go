package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gempir/gempbot/pkg/auth"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/log"
	nickHelix "github.com/nicklaw5/helix"
)

type ApiClient struct {
	cfg        *config.Config
	httpClient *http.Client
	baseUrl    string
}

func NewApiClient(cfg *config.Config) *ApiClient {
	return &ApiClient{
		cfg:        cfg,
		httpClient: &http.Client{},
		baseUrl:    cfg.ApiBaseUrl,
	}
}

func (c *ApiClient) CreatePrediction(userID string, prediction *nickHelix.CreatePredictionParams) (*http.Response, error) {
	marshalled, err := json.Marshal(prediction)
	if err != nil {
		return nil, err
	}

	log.Infof("[POST] %s/api/prediction", c.baseUrl)
	req, err := http.NewRequest(http.MethodPost, c.baseUrl+"/api/prediction", bytes.NewBuffer(marshalled))
	req.Header.Set("Authorization", "Bearer "+auth.CreateApiToken(c.cfg.Secret, userID))
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return resp, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("%s", body)
	}

	return resp, err
}
