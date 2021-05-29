package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/slice"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/labstack/echo/v4"
)

type UserConfig struct {
	Editors   []string
	Protected Protected
}

type Protected struct {
	EditorFor     []string
	CurrentUserID string
}

func createDefaultUserConfig() UserConfig {
	return UserConfig{
		Editors: []string{},
		Protected: Protected{
			EditorFor:     []string{},
			CurrentUserID: "",
		},
	}
}

func (c *UserConfig) getEditorDifference(newEditors []string) (removed []string, added []string) {
	return slice.Diff(c.Editors, newEditors)
}

func (s *Server) handleUserConfig(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet {
		userConfig := s.getUserConfig(auth.Data.UserID)

		if c.QueryParam("managing") != "" {
			ownerUserID, err := s.checkEditor(c, userConfig)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			editorFor := userConfig.Protected.EditorFor
			userConfig = s.getUserConfig(ownerUserID)

			userConfig.Editors = []string{}
			userConfig.Protected.EditorFor = editorFor
		}

		userConfig = s.convertUserConfig(userConfig, true)

		return c.JSON(http.StatusOK, userConfig)

	} else if c.Request().Method == http.MethodPost {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Errorf("Failed reading update body: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failure saving body "+err.Error())
		}

		var newConfig UserConfig
		if err := json.Unmarshal(body, &newConfig); err != nil {
			log.Errorf("Failed unmarshalling userConfig: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failure unmarshalling config "+err.Error())
		}

		if c.QueryParam("managing") != "" {
			return echo.NewHTTPError(http.StatusForbidden, "editors are not allowed to edit userConfig yet")
		}

		err = s.processConfig(auth.Data.UserID, newConfig, c)
		if err != nil {
			log.Errorf("failed processing config: %s", err)
			return echo.NewHTTPError(http.StatusBadRequest, "failed processing config: "+err.Error())
		}

		return c.JSON(http.StatusOK, nil)
	}

	return nil
}

func (s *Server) getUserConfig(userID string) UserConfig {
	editors := s.db.GetEditors(userID)

	uCfg := createDefaultUserConfig()
	uCfg.Protected.CurrentUserID = userID

	for _, editor := range editors {
		if editor.OwnerTwitchID == userID {
			uCfg.Editors = append(uCfg.Editors, editor.EditorTwitchID)
		}
		if editor.EditorTwitchID == userID {
			uCfg.Protected.EditorFor = append(uCfg.Protected.EditorFor, editor.OwnerTwitchID)
		}
	}

	return uCfg
}

func (s *Server) convertUserConfig(uCfg UserConfig, toNames bool) UserConfig {
	all := uCfg.Editors
	all = append(all, uCfg.Protected.EditorFor...)

	var err error
	var userData map[string]helix.UserData
	if toNames {
		userData, err = s.helixClient.GetUsersByUserIds(all)
	} else {
		userData, err = s.helixClient.GetUsersByUsernames(all)
	}
	if err != nil {
		log.Errorf("Failed to get editors %s", err)
		return UserConfig{}
	}

	editors := []string{}
	for _, editor := range uCfg.Editors {
		data, ok := userData[editor]
		if !ok {
			continue
		}

		if toNames {
			editors = append(editors, data.Login)
		} else {
			editors = append(editors, data.ID)
		}
	}
	uCfg.Editors = editors

	editorFor := []string{}
	for _, editor := range uCfg.Protected.EditorFor {
		data, ok := userData[editor]
		if !ok {
			continue
		}

		if toNames {
			editorFor = append(editorFor, data.Login)
		} else {
			editorFor = append(editorFor, data.ID)
		}
	}
	uCfg.Protected.EditorFor = editorFor

	return uCfg
}

func (s *Server) checkEditor(c echo.Context, userConfig UserConfig) (string, error) {
	managing := c.QueryParam("managing")

	if managing == "" {
		return "", nil
	}

	userData, err := s.helixClient.GetUsersByUsernames([]string{managing})
	if err != nil || len(userData) == 0 {
		return "", errors.New("can't find editor")
	}

	isEditor := false
	for _, editor := range userConfig.Protected.EditorFor {
		if editor == userData[managing].ID {
			isEditor = true
		}
	}

	if !isEditor {
		return "", errors.New("User is not editor")
	}

	return userData[managing].ID, nil
}

func (s *Server) checkIsEditor(editorUserID string, ownerUserID string) error {
	if editorUserID == ownerUserID {
		return nil
	}

	userConfig := s.getUserConfig(ownerUserID)

	for _, editor := range userConfig.Editors {
		if editor == editorUserID {
			return nil
		}
	}

	return echo.NewHTTPError(http.StatusForbidden, "user is not editor")
}

func (s *Server) processConfig(userID string, newConfig UserConfig, c echo.Context) error {
	newUserIDConfig := s.convertUserConfig(newConfig, false)
	oldConfig := s.getUserConfig(userID)
	added, removed := oldConfig.getEditorDifference(newUserIDConfig.Editors)

	s.db.AddEditors(userID, added)
	s.db.RemoveEditors(userID, removed)

	s.subscribeChannelPoints(userID)

	return nil
}

type MigrateUserConfig struct {
	Editors []string
	Rewards MigrateRewards
}

type MigrateRewards struct {
	*MigrateBttvReward `json:"Bttv"`
}

type MigrateBttvReward struct {
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

type MigrateUserAcessTokenData struct {
	AccessToken  string
	RefreshToken string
	Scope        string
}

func (s *Server) migrateData() {
	configs, err := s.store.Client.HGetAll("userConfig").Result()
	if err != nil {
		log.Info(err)
	}

	for userID, dataString := range configs {
		var uCfg MigrateUserConfig
		err := json.Unmarshal([]byte(dataString), &uCfg)
		if err != nil {
			log.Error(err)
		}

		s.db.AddEditors(userID, uCfg.Editors)

		if uCfg.Rewards.MigrateBttvReward != nil && uCfg.Rewards.MigrateBttvReward.ID != "" {
			log.Infof("bttv reward created for %s", userID)
			rew := uCfg.Rewards.MigrateBttvReward
			err := s.db.SaveReward(store.ChannelPointReward{
				OwnerTwitchID:                     userID,
				Type:                              TYPE_BTTV,
				Title:                             rew.Title,
				Prompt:                            rew.Prompt,
				Cost:                              rew.Cost,
				BackgroundColor:                   rew.Backgroundcolor,
				IsMaxPerStreamEnabled:             rew.IsMaxPerStreamEnabled,
				MaxPerStream:                      rew.MaxPerStream,
				IsUserInputRequired:               rew.IsUserInputRequired,
				IsMaxPerUserPerStreamEnabled:      rew.IsMaxPerUserPerStreamEnabled,
				MaxPerUserPerStream:               rew.MaxPerUserPerStream,
				IsGlobalCooldownEnabled:           rew.IsGlobalCooldownEnabled,
				GlobalCooldownSeconds:             rew.GlobalCooldownSeconds,
				ShouldRedemptionsSkipRequestQueue: rew.ShouldRedemptionsSkipRequestQueue,
				Enabled:                           rew.Enabled,
				RewardID:                          rew.ID,
			})
			if err != nil {
				log.Error(err)
			}
		}
	}

	tokens, err := s.store.Client.HGetAll("userAccessTokensData").Result()
	if err != nil {
		log.Info(err)
	}

	for userID, dataString := range tokens {
		var tokenData MigrateUserAcessTokenData
		err := json.Unmarshal([]byte(dataString), &tokenData)
		if err != nil {
			log.Error(err)
		}

		if tokenData.AccessToken != "" {
			log.Infof("token created for %s", userID)
			err := s.db.SaveUserAccessToken(userID, tokenData.AccessToken, tokenData.RefreshToken, tokenData.Scope)
			if err != nil {
				log.Error(err)
			}
		}
	}

	bttvs, err := s.store.Client.HGetAll("bttv_emote").Result()
	if err != nil {
		log.Info(err)
	}

	for userID, emoteID := range bttvs {
		log.Infof("bttv emote created for %s %s", userID, emoteID)
		s.db.CreateEmoteAdd(userID, emoteID)
	}
}
