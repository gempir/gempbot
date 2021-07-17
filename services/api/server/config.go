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
	BotJoin     bool
	Permissions map[string]Permission
	Protected   Protected
}

func (u *UserConfig) isEditor(userID string) bool {
	val, ok := u.Permissions[userID]

	return ok && val.Editor
}

func (u *UserConfig) isEditorFor(userID string) bool {
	return slice.Contains(u.Protected.EditorFor, userID)
}

type Protected struct {
	EditorFor     []string
	CurrentUserID string
}

type Permission struct {
	Editor     bool
	Prediction bool
}

func createDefaultUserConfig() UserConfig {
	return UserConfig{
		BotJoin:     false,
		Permissions: map[string]Permission{},
		Protected: Protected{
			EditorFor:     []string{},
			CurrentUserID: "",
		},
	}
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
				return err
			}

			editorFor := userConfig.Protected.EditorFor
			userConfig = s.getUserConfig(ownerUserID)

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

		err = s.processConfig(auth.Data.UserID, auth.Data.Login, newConfig, c.QueryParam("managing"))
		if err != nil {
			log.Errorf("failed processing config: %s", err)
			return echo.NewHTTPError(http.StatusBadRequest, "failed processing config: "+err.Error())
		}

		return c.JSON(http.StatusOK, nil)
	}

	return nil
}

func (s *Server) getUserConfig(userID string) UserConfig {
	uCfg := createDefaultUserConfig()

	botConfig, err := s.db.GetBotConfig(userID)
	if err != nil {
		uCfg.BotJoin = false
	} else {
		uCfg.BotJoin = botConfig.JoinBot
	}

	uCfg.Protected.CurrentUserID = userID

	perms := s.db.GetChannelPermissions(userID)
	for _, perm := range perms {
		uCfg.Permissions[perm.TwitchID] = Permission{perm.Editor, perm.Prediction}
	}

	for _, perm := range s.db.GetUserPermissions(userID) {
		uCfg.Protected.EditorFor = append(uCfg.Protected.EditorFor, perm.ChannelTwitchId)
	}

	return uCfg
}

func (s *Server) convertUserConfig(uCfg UserConfig, toNames bool) UserConfig {
	all := uCfg.Protected.EditorFor

	for user := range uCfg.Permissions {
		all = append(all, user)
	}

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

	perms := map[string]Permission{}
	for user, perm := range uCfg.Permissions {
		data, ok := userData[user]
		if !ok {
			continue
		}

		var newUser string

		if toNames {
			newUser = data.Login
		} else {
			newUser = data.ID
		}

		perms[newUser] = perm
	}
	uCfg.Permissions = perms

	return uCfg
}

func (s *Server) checkEditor(c echo.Context, userConfig UserConfig) (string, error) {
	managing := c.QueryParam("managing")

	userData, err := s.helixClient.GetUsersByUsernames([]string{managing})
	if err != nil || len(userData) == 0 {
		return "", echo.NewHTTPError(http.StatusForbidden, "could not find managing")
	}

	if !userConfig.isEditorFor(userData[managing].ID) {
		return "", echo.NewHTTPError(http.StatusForbidden, "user is not editor")
	}

	return userData[managing].ID, nil
}

func (s *Server) checkIsEditor(editorUserID string, ownerUserID string) error {
	if editorUserID == ownerUserID {
		return nil
	}

	userConfig := s.getUserConfig(ownerUserID)

	if userConfig.isEditor(editorUserID) {
		return nil
	}

	return echo.NewHTTPError(http.StatusForbidden, "user is not editor")
}

func (s *Server) processConfig(userID string, login string, newConfig UserConfig, managing string) error {
	isManaging := managing != ""
	ownerUserID := userID
	ownerLogin := login
	newUserIDConfig := s.convertUserConfig(newConfig, false)

	if isManaging {
		uData, err := s.helixClient.GetUserByUsername(managing)
		if err != nil {
			return err
		}
		ownerUserID = uData.ID
		ownerLogin = uData.Login
		oldConfig := s.getUserConfig(uData.ID)

		if !oldConfig.isEditor(userID) {
			return errors.New("not an editor")
		}
	}

	err := s.db.SaveBotConfig(store.BotConfig{OwnerTwitchID: ownerUserID, JoinBot: newUserIDConfig.BotJoin})
	if err != nil {
		log.Error(err)
	}
	if newConfig.BotJoin {
		s.store.PublishIngesterMessage(store.IngesterMsgJoin, ownerLogin)
	} else {
		s.store.PublishIngesterMessage(store.IngesterMsgPart, ownerLogin)
	}

	previousPerms := s.db.GetChannelPermissions(ownerUserID)
	previousPermIds := []string{}
	for _, perm := range previousPerms {
		previousPermIds = append(previousPermIds, perm.TwitchID)
	}
	newPermIds := []string{}
	for user := range newUserIDConfig.Permissions {
		newPermIds = append(newPermIds, user)
	}
	_, deleted := slice.Diff(previousPermIds, newPermIds)

	for _, deletedUserID := range deleted {
		s.db.DeletePermission(ownerUserID, deletedUserID)
	}

	for permissionUserID, perm := range newUserIDConfig.Permissions {
		newPerms := map[string]interface{}{"ChannelTwitchId": ownerUserID, "TwitchID": permissionUserID, "Prediction": perm.Prediction}
		if !isManaging {
			newPerms["Editor"] = perm.Editor
		}

		err := s.db.SavePermission(newPerms)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
