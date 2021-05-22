package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

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

	cfg, err, isNew := s.getUserConfig(cfgUserID)
	if err != nil || isNew {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("no config found %s", err))
	}

	if cfg.Rewards.BttvReward.IsDefault || cfg.Rewards.BttvReward.ID != rewardID {
		return echo.NewHTTPError(http.StatusBadRequest, "rewardId not found in config")
	}

	token, err := s.getUserAccessToken(cfgUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "no accessToken to edit reward")
	}

	cfg.Rewards.BttvReward = createDefaultBttvReward()

	err = s.saveConfig(cfgUserID, cfg)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.helixUserClient.DeleteReward(cfgUserID, token.AccessToken, rewardID)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
