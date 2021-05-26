package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gempir/bitraft/pkg/store"
	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type UserConfig struct {
	Editors       []string
	Protected     Protected
	Rewards       []Reward
	CurrentUserID string
}

type Protected struct {
	EditorFor []string
}

func createDefaultUserConfig() UserConfig {
	return UserConfig{
		Editors: []string{},
		Protected: Protected{
			EditorFor: []string{},
		},
		Rewards: []Reward{},
	}
}

func (s *Server) handleUserConfig(c echo.Context) error {
	auth, _, err := s.authenticate(c)
	if err != nil {
		return err
	}

	if c.Request().Method == http.MethodGet {
		userConfig, err, _ := s.getUserConfig(auth.Data.UserID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "can't recover config"+err.Error())
		}

		managing := c.QueryParam("managing")
		if managing != "" {
			ownerUserID, err := s.checkEditor(c, userConfig)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			managingUserConfig, err, _ := s.getUserConfig(ownerUserID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "can't recover config"+err.Error())
			}
			newEditorForNames := []string{}

			var editors []store.Editor
			s.db.Where("editor_twitch_id = ?", auth.Data.UserID).Find(&editors)

			editorForIds := []string{}
			for _, editor := range editors {
				editorForIds = append(editorForIds, editor.OwnerTwitchID)
			}

			userData, err := s.helixClient.GetUsersByUserIds(editorForIds)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "can't resolve editorFor in config "+err.Error())
			}
			for _, user := range userData {
				newEditorForNames = append(newEditorForNames, user.Login)
			}

			managingUserConfig.Editors = []string{}
			managingUserConfig.Protected.EditorFor = newEditorForNames
			managingUserConfig.CurrentUserID = ownerUserID

			return c.JSON(http.StatusOK, managingUserConfig)
		} else {
			newEditorNames := []string{}

			var editors []store.Editor
			s.db.Where("owner_twitch_id = ?", auth.Data.UserID).Find(&editors)

			editorIds := []string{}
			for _, editor := range editors {
				editorIds = append(editorIds, editor.EditorTwitchID)
			}

			userData, err := s.helixClient.GetUsersByUserIds(editorIds)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "can't resolve editorFor in config "+err.Error())
			}
			for _, user := range userData {
				newEditorNames = append(newEditorNames, user.Login)
			}

			userConfig.Editors = newEditorNames

			newEditorForNames := []string{}

			userData, err = s.helixClient.GetUsersByUserIds(userConfig.Protected.EditorFor)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "can't resolve editorFor in config "+err.Error())
			}
			for _, user := range userData {
				newEditorForNames = append(newEditorForNames, user.Login)
			}

			userConfig.CurrentUserID = auth.Data.UserID

			userConfig.Protected.EditorFor = newEditorForNames
			return c.JSON(http.StatusOK, userConfig)
		}

	} else if c.Request().Method == http.MethodPost {
		body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Errorf("Failed reading update body: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failure saving body "+err.Error())
		}

		_, err = s.processConfig(auth.Data.UserID, body, c)
		if err != nil {
			log.Errorf("failed processing config: %s", err)
			return echo.NewHTTPError(http.StatusBadRequest, "failed processing config: "+err.Error())
		}

		return c.JSON(http.StatusOK, nil)
	} else if c.Request().Method == http.MethodDelete {
		_, err := s.store.Client.HDel("userConfig", auth.Data.UserID).Result()
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed deleting: "+err.Error())
		}

		err = s.unsubscribeChannelPoints(auth.Data.UserID, "userDeleted")
		if err != nil {
			log.Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to unsubscribe: "+err.Error())
		}

		return c.JSON(http.StatusOK, nil)
	}

	return nil
}

func (s *Server) getUserConfig(userID string) (UserConfig, error, bool) {
	val, err := s.store.Client.HGet("userConfig", userID).Result()
	if err == redis.Nil || val == "" {
		return createDefaultUserConfig(), nil, true
	}
	if err != nil {
		return createDefaultUserConfig(), err, true
	}

	var userConfig UserConfig
	if err := json.Unmarshal([]byte(val), &userConfig); err != nil {
		return createDefaultUserConfig(), errors.New("can't find config"), true
	}

	return userConfig, nil, false
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

	userConfig, err, _ := s.getUserConfig(ownerUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "no config found for owner")
	}

	for _, editor := range userConfig.Editors {
		if editor == editorUserID {
			return nil
		}
	}

	return echo.NewHTTPError(http.StatusForbidden, "user is not editor")
}

func (s *Server) processConfig(userID string, body []byte, c echo.Context) (UserConfig, error) {
	oldConfig, err, isNew := s.getUserConfig(userID)
	if err != nil {
		return UserConfig{}, err
	}
	editorConfig := oldConfig

	ownerUserID, err := s.checkEditor(c, oldConfig)
	if err != nil {
		return UserConfig{}, err
	}

	if ownerUserID != "" {
		oldConfig, err, isNew = s.getUserConfig(ownerUserID)
		if err != nil {
			return UserConfig{}, err
		}
	}

	var newConfig UserConfig
	if err := json.Unmarshal(body, &newConfig); err != nil {
		return UserConfig{}, err
	}

	protected := oldConfig.Protected
	if protected.EditorFor == nil {
		protected.EditorFor = []string{}
	}

	configToSave := UserConfig{
		Editors:   oldConfig.Editors,
		Protected: protected,
	}

	if ownerUserID == "" {
		newEditorIds := []string{}

		userData, err := s.helixClient.GetUsersByUsernames(newConfig.Editors)
		if err != nil {
			return UserConfig{}, err
		}
		if len(newConfig.Editors) != len(userData) {
			return UserConfig{}, errors.New("Failed to find all editors")
		}

		for _, user := range userData {
			if user.ID == userID {
				return UserConfig{}, errors.New("You can't be your own editor")
			}

			newEditorIds = append(newEditorIds, user.ID)
		}

		configToSave.Editors = newEditorIds

		for _, editor := range oldConfig.Editors {
			err := s.removeEditorFor(editor, userID)
			if err != nil {
				return UserConfig{}, err
			}
		}

		for _, user := range userData {
			err := s.addEditorFor(user.ID, userID)
			if err != nil {
				return UserConfig{}, err
			}
		}
	}

	saveTarget := userID
	if ownerUserID != "" {
		saveTarget = ownerUserID
	}

	for _, reward := range newConfig.Rewards {
		if reward.GetType() == TYPE_BTTV {
			oldRewardId := ""
			for _, oldReward := range oldConfig.Rewards {
				if oldReward.GetType() == TYPE_BTTV {
					oldRewardId = oldReward.GetConfig().ID
				}
			}

			_, err := s.createOrUpdateChannelPointReward(saveTarget, reward, oldRewardId)
			if err != nil {
				return UserConfig{}, err
			}

		}
	}

	err = s.saveConfig(saveTarget, configToSave)
	if err != nil {
		return UserConfig{}, err
	}

	if isNew {
		log.Infof("Created new config for: %s, subscribing webhooks", userID)
		s.subscribeChannelPoints(userID)
	}

	configToSave.Editors = []string{}
	if ownerUserID == "" {
		configToSave.Editors = newConfig.Editors
	}

	configToSave.Protected = editorConfig.Protected

	newEditorForNames := []string{}

	userData, err := s.helixClient.GetUsersByUserIds(editorConfig.Protected.EditorFor)
	if err != nil {
		return UserConfig{}, errors.New("can't resolve editorFor in config " + err.Error())
	}
	for _, user := range userData {
		newEditorForNames = append(newEditorForNames, user.Login)
	}

	configToSave.Protected.EditorFor = newEditorForNames

	return configToSave, nil
}

func (s *Server) saveConfig(userID string, userConfig UserConfig) error {
	js, err := json.Marshal(userConfig)
	if err != nil {
		return err
	}

	_, err = s.store.Client.HSet("userConfig", userID, js).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) addEditorFor(editorID, userID string) error {
	isNew := false

	val, err := s.store.Client.HGet("userConfig", editorID).Result()
	if err == redis.Nil {
		isNew = true
	} else if err != nil {
		return err
	}

	var userConfig UserConfig
	if !isNew {
		if err := json.Unmarshal([]byte(val), &userConfig); err != nil {
			return err
		}
	} else {
		userConfig = createDefaultUserConfig()
	}

	userConfig.Protected.EditorFor = append(userConfig.Protected.EditorFor, userID)

	js, err := json.Marshal(userConfig)
	if err != nil {
		return err
	}

	log.Infof("New Editor %s for user %s", editorID, userID)
	_, err = s.store.Client.HSet("userConfig", editorID, js).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) removeEditorFor(editorID, userID string) error {
	isNew := true

	val, err := s.store.Client.HGet("userConfig", editorID).Result()
	if err == redis.Nil {
		isNew = true
	} else if err != nil {
		return err
	}

	var userConfig UserConfig
	if !isNew {
		if err := json.Unmarshal([]byte(val), &userConfig); err != nil {
			return err
		}
	} else {
		userConfig = createDefaultUserConfig()
	}

	newEditorFor := []string{}

	for _, editor := range userConfig.Protected.EditorFor {
		if userID != editor {
			newEditorFor = append(newEditorFor, editor)
		}
	}

	userConfig.Protected.EditorFor = newEditorFor

	js, err := json.Marshal(userConfig)
	if err != nil {
		return err
	}

	log.Infof("New Editor %s for user %s", editorID, userID)
	_, err = s.store.Client.HSet("userConfig", editorID, js).Result()
	if err != nil {
		return err
	}

	return nil
}
