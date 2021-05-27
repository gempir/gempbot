package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/slice"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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

func (c *UserConfig) isEditorFor(user string) bool {
	for _, editor := range c.Protected.EditorFor {
		if editor == user {
			return true
		}
	}

	return false
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

			userConfig = s.getUserConfig(ownerUserID)
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
	var editors []store.Editor
	s.db.Where("owner_twitch_id = ? OR editor_twitch_id = ?", userID, userID).Find(&editors)

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

func (s *Server) addEditors(ownerId string, userIds []string) {
	if len(userIds) == 0 {
		return
	}

	var editors []store.Editor
	for _, id := range userIds {
		editors = append(editors, store.Editor{OwnerTwitchID: ownerId, EditorTwitchID: id})
	}

	s.db.Create(&editors)
}

func (s *Server) removeEditors(ownerId string, userIds []string) {
	if len(userIds) == 0 {
		return
	}

	s.db.Delete(store.Editor{}, "editor_twitch_id IN (?) AND owner_user_id = ?", strings.Join(userIds, ","), ownerId)
}

func (s *Server) processConfig(userID string, newConfig UserConfig, c echo.Context) error {
	newUserIDConfig := s.convertUserConfig(newConfig, false)
	oldConfig := s.getUserConfig(userID)
	added, removed := oldConfig.getEditorDifference(newUserIDConfig.Editors)

	s.addEditors(userID, added)
	s.removeEditors(userID, removed)

	// for _, reward := range newConfig.Rewards {
	// 	if reward.GetType() == TYPE_BTTV {
	// 		oldRewardId := ""
	// 		for _, oldReward := range oldConfig.Rewards {
	// 			if oldReward.GetType() == TYPE_BTTV {
	// 				oldRewardId = oldReward.GetConfig().ID
	// 			}
	// 		}

	// 		_, err := s.createOrUpdateChannelPointReward(saveTarget, reward, oldRewardId)
	// 		if err != nil {
	// 			return UserConfig{}, err
	// 		}

	// 	}
	// }

	// if isNew {
	// 	log.Infof("Created new config for: %s, subscribing webhooks", userID)
	// 	s.subscribeChannelPoints(userID)
	// }

	return nil
}
