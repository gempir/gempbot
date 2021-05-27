package server

import (
	"encoding/json"
	"net/http"

	"github.com/gempir/bitraft/pkg/store"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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

	cfgUserID := c.Param("userID")
	rewardID := c.Param("rewardID")

	err = s.checkIsEditor(auth.Data.UserID, cfgUserID)
	if err != nil {
		return err
	}

	// // cfg, err, isNew := s.getUserConfig(cfgUserID)
	// // if err != nil || isNew {
	// // 	return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("no config found %s", err))
	// }

	// updatedRewards := []Reward{}

	// for _, reward := range cfg.Rewards {
	// 	if reward.GetConfig().ID != rewardID {
	// 		updatedRewards = append(updatedRewards, reward)
	// 	}
	// }

	// cfg.Rewards = updatedRewards

	token, err := s.getUserAccessToken(cfgUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no accessToken to edit reward")
	}

	// err = s.saveConfig(cfgUserID, cfg)
	// if err != nil {
	// 	log.Error(err)
	// 	return err
	// }

	err = s.helixUserClient.DeleteReward(cfgUserID, token.AccessToken, rewardID)
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

	if c.QueryParam("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.QueryParam("userID"))
		if err != nil {
			return err
		}
	}

	rewards := s.db.GetChannelPointRewards(c.QueryParam("userID"))

	return c.JSON(http.StatusOK, rewards)
}

func (s *Server) handleRewardSingleRead(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.QueryParam("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.QueryParam("userID"))
		if err != nil {
			return err
		}
	}

	reward := s.db.GetChannelPointReward(c.QueryParam("userID"), c.QueryParam("type"))

	return c.JSON(http.StatusOK, reward)
}

func (s *Server) handleRewardCreateOrUpdate(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.QueryParam("userID") != auth.Data.UserID {
		err := s.checkIsEditor(auth.Data.UserID, c.QueryParam("userID"))
		if err != nil {
			return err
		}
	}

	var newReward store.ChannelPointReward
	if err := json.NewDecoder(c.Request().Body).Decode(&newReward); err != nil {
		log.Errorf("Failed unmarshalling reward: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failure unmarshalling reward")
	}

	newReward = setRewardDefaults(newReward, c.QueryParam("userID"))

	err = s.db.SaveReward(newReward)
	if err != nil {
		log.Errorf("Failed saving reward: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failure saving reward")
	}

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
