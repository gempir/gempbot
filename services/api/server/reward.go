package server

import (
	"encoding/json"
	"net/http"

	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/labstack/echo/v4"
)

const (
	TYPE_BTTV = "bttv"
)

type Reward interface {
	GetType() string
	GetConfig() TwitchRewardConfig
}

type TwitchRewardConfig struct {
	Title                             string `json:"title"`
	Prompt                            string `json:"prompt"`
	Cost                              int    `json:"cost"`
	Backgroundcolor                   string `json:"backgroundColor"`
	IsMaxPerStreamEnabled             bool   `json:"isMaxPerStreamEnabled"`
	MaxPerStream                      int    `json:"maxPerStream"`
	IsUserInputRequired               bool   `json:"isUserInputRequired"`
	IsMaxPerUserPerStreamEnabled      bool   `json:"isMaxPerUserPerStreamEnabled"`
	MaxPerUserPerStream               int    `json:"maxPerUserPerStream"`
	IsGlobalCooldownEnabled           bool   `json:"isGlobalCooldownEnabled"`
	GlobalCooldownSeconds             int    `json:"globalCooldownSeconds"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"shouldRedemptionsSkipRequestQueue"`
	Enabled                           bool   `json:"enabled"`
	ID                                string
}

type BttvReward struct {
	TwitchRewardConfig
}

func (r *BttvReward) GetType() string {
	return TYPE_BTTV
}

func MarshallReward(reward Reward) string {
	js, err := json.Marshal(reward)
	if err != nil {
		log.Infof("failed to marshal BttvReward %s", err)
		return ""
	}

	return string(js)
}

func (s *Server) handleRewardDeletion(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Param("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.Param("userID"))
		if err != nil {
			return err
		}
	}

	reward, err := s.db.GetChannelPointReward(c.Param("userID"), c.Param("type"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	token, err := s.db.GetUserAccessToken(c.Param("userID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no accessToken to edit reward")
	}

	s.db.DeleteChannelPointReward(c.Param("userID"), c.Param("type"))

	err = s.helixUserClient.DeleteReward(c.Param("userID"), token.AccessToken, reward.RewardID)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) handleRewardRead(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Param("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.Param("userID"))
		if err != nil {
			return err
		}
	}

	rewards := s.db.GetChannelPointRewards(c.Param("userID"))

	return c.JSON(http.StatusOK, rewards)
}

func (s *Server) handleRewardSingleRead(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Param("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.Param("userID"))
		if err != nil {
			return err
		}
	}

	reward, err := s.db.GetChannelPointReward(c.Param("userID"), c.Param("type"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, reward)
}

func (s *Server) handleRewardCreateOrUpdate(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Param("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.Param("userID"))
		if err != nil {
			return err
		}
	}

	var newReward store.ChannelPointReward
	if err := json.NewDecoder(c.Request().Body).Decode(&newReward); err != nil {
		log.Errorf("Failed unmarshalling reward: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failure unmarshalling reward")
	}

	newReward = setRewardDefaults(newReward, c.Param("userID"))

	rewardID := ""
	reward, err := s.db.GetChannelPointReward(c.Param("userID"), TYPE_BTTV)
	if err == nil {
		rewardID = reward.RewardID
	}

	newReward.RewardID, err = s.createOrUpdateChannelPointReward(c.Param("userID"), newReward, rewardID)
	if err != nil {
		log.Errorf("Failed saving reward to twitch: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failure saving reward to twitch")
	}

	err = s.db.SaveReward(newReward)
	if err != nil {
		log.Errorf("Failed saving reward: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failure saving reward")
	}

	s.subscribeChannelPoints(c.Param("userID"))

	return nil
}

func setRewardDefaults(reward store.ChannelPointReward, userID string) store.ChannelPointReward {
	reward.OwnerTwitchID = userID
	reward.IsUserInputRequired = true
	reward.ShouldRedemptionsSkipRequestQueue = false

	if reward.Type == TYPE_BTTV {
		reward.Prompt = bttvPrompt
	}

	return reward
}

func (s *Server) createOrUpdateChannelPointReward(userID string, request store.ChannelPointReward, rewardID string) (string, error) {
	token, err := s.db.GetUserAccessToken(userID)
	if err != nil {
		return "", err
	}

	req := helix.CreateCustomRewardRequest{
		Title:                             request.Title,
		Prompt:                            bttvPrompt,
		Cost:                              request.Cost,
		IsEnabled:                         request.Enabled,
		BackgroundColor:                   request.BackgroundColor,
		IsUserInputRequired:               true,
		ShouldRedemptionsSkipRequestQueue: false,
		IsMaxPerStreamEnabled:             false,
		IsMaxPerUserPerStreamEnabled:      false,
		IsGlobalCooldownEnabled:           false,
	}

	if request.MaxPerStream != 0 {
		req.IsMaxPerStreamEnabled = true
		req.MaxPerStream = request.MaxPerStream
	}

	if request.MaxPerUserPerStream != 0 {
		req.IsMaxPerUserPerStreamEnabled = true
		req.MaxPerUserPerStream = request.MaxPerUserPerStream
	}

	if request.GlobalCooldownSeconds != 0 {
		req.IsGlobalCooldownEnabled = true
		req.GlobalCoolDownSeconds = request.GlobalCooldownSeconds
	}

	resp, err := s.helixUserClient.CreateOrUpdateReward(userID, token.AccessToken, req, rewardID)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}
