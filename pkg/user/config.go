package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gempir/gempbot/pkg/api"
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/slice"
	"github.com/gempir/gempbot/pkg/store"
)

type UserAdmin struct {
	cfg         *config.Config
	db          *store.Database
	helixClient *helix.Client
	chatClient  *chat.ChatClient
}

func NewUserAdmin(cfg *config.Config, db *store.Database, helixClient *helix.Client, chatClient *chat.ChatClient) *UserAdmin {
	return &UserAdmin{cfg, db, helixClient, chatClient}
}

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

func (u *UserAdmin) GetUserConfig(userID string) UserConfig {
	uCfg := createDefaultUserConfig()

	botConfig, err := u.db.GetBotConfig(userID)
	if err != nil {
		uCfg.BotJoin = false
	} else {
		uCfg.BotJoin = botConfig.JoinBot
	}

	uCfg.Protected.CurrentUserID = userID

	perms := u.db.GetChannelPermissions(userID)
	for _, perm := range perms {
		uCfg.Permissions[perm.TwitchID] = Permission{perm.Editor, perm.Prediction}
	}

	for _, perm := range u.db.GetUserPermissions(userID) {
		uCfg.Protected.EditorFor = append(uCfg.Protected.EditorFor, perm.ChannelTwitchId)
	}

	return uCfg
}

func (u *UserAdmin) ConvertUserConfig(uCfg UserConfig, toNames bool) (UserConfig, api.Error) {
	all := map[string]string{}

	for _, user := range uCfg.Protected.EditorFor {
		all[user] = user
	}

	for user := range uCfg.Permissions {
		all[user] = user
	}

	allSlice := slice.MapToSlice(all)

	var err error
	var userData map[string]helix.UserData
	if toNames {
		userData, err = u.helixClient.GetUsersByUserIds(allSlice)
	} else {
		userData, err = u.helixClient.GetUsersByUsernames(allSlice)
	}
	if err != nil || len(userData) != len(all) {
		log.Errorf("Failed to get users %s", err)
		return UserConfig{}, api.NewApiError(http.StatusBadRequest, errors.New("failed to get users"))
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

	return uCfg, nil
}

func (u *UserAdmin) CheckEditor(r *http.Request, userConfig UserConfig) (string, api.Error) {
	managing := r.URL.Query().Get("managing")

	userData, err := u.helixClient.GetUsersByUsernames([]string{managing})
	if err != nil || len(userData) == 0 {
		return "", api.NewApiError(http.StatusForbidden, fmt.Errorf("could not find managing"))
	}

	if !userConfig.isEditorFor(userData[managing].ID) {
		return "", api.NewApiError(http.StatusForbidden, fmt.Errorf("user is not editor"))
	}

	return userData[managing].ID, nil
}

func (u *UserAdmin) ProcessConfig(ctx context.Context, userID string, login string, newConfig UserConfig, managing string) api.Error {
	isManaging := managing != ""
	ownerUserID := userID
	ownerLogin := login
	newUserIDConfig, err := u.ConvertUserConfig(newConfig, false)
	if err != nil {
		return err
	}

	if isManaging {
		uData, err := u.helixClient.GetUserByUsername(managing)
		if err != nil {
			return api.NewApiError(http.StatusBadRequest, fmt.Errorf("could not find managing"))
		}
		ownerUserID = uData.ID
		ownerLogin = uData.Login
		oldConfig := u.GetUserConfig(uData.ID)

		if !oldConfig.isEditor(userID) {
			return api.NewApiError(http.StatusForbidden, fmt.Errorf("user is not editor"))
		}
	}

	dbErr := u.db.SaveBotConfig(ctx, store.BotConfig{OwnerTwitchID: ownerUserID, JoinBot: newUserIDConfig.BotJoin})
	if dbErr != nil {
		log.Error(dbErr)
		return api.NewApiError(http.StatusInternalServerError, fmt.Errorf("failed to save bot config"))
	}
	if newConfig.BotJoin {
		u.chatClient.JoinBot(ownerLogin)
	} else {
		u.chatClient.PartBot(ownerLogin)
	}

	previousPerms := u.db.GetChannelPermissions(ownerUserID)
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
		u.db.DeletePermission(ownerUserID, deletedUserID)
	}

	for permissionUserID, perm := range newUserIDConfig.Permissions {
		newPerms := store.Permission{ChannelTwitchId: ownerUserID, TwitchID: permissionUserID, Prediction: perm.Prediction, Editor: perm.Editor}

		err := u.db.SavePermission(newPerms)
		if err != nil {
			log.Error(err)
		}
	}

	return nil
}
