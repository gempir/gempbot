package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gempir/gempbot/internal/api"
	"github.com/gempir/gempbot/internal/chat"
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/slice"
	"github.com/gempir/gempbot/internal/store"
)

type UserAdmin struct {
	cfg         *config.Config
	db          *store.Database
	helixClient *helixclient.Client
	chatClient  *chat.ChatClient
}

func NewUserAdmin(cfg *config.Config, db *store.Database, helixClient *helixclient.Client, chatClient *chat.ChatClient) *UserAdmin {
	return &UserAdmin{cfg, db, helixClient, chatClient}
}

type UserConfig struct {
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
		Permissions: map[string]Permission{},
		Protected: Protected{
			EditorFor:     []string{},
			CurrentUserID: "",
		},
	}
}

func (u *UserAdmin) GetUserConfig(userID string) UserConfig {
	uCfg := createDefaultUserConfig()

	uCfg.Protected.CurrentUserID = userID

	perms := u.db.GetChannelPermissions(userID)
	for _, perm := range perms {
		uCfg.Permissions[perm.TwitchID] = Permission{perm.Editor, perm.Prediction}
	}

	for _, perm := range u.db.GetUserPermissions(userID) {
		if perm.Editor {
			uCfg.Protected.EditorFor = append(uCfg.Protected.EditorFor, perm.ChannelTwitchId)
		}
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
	var userData map[string]helixclient.UserData
	if toNames {
		userData, err = u.helixClient.GetUsersByUserIds(allSlice)
	} else {
		userData, err = u.helixClient.GetUsersByUsernames(allSlice)
	}
	if err != nil || len(userData) != len(all) {
		log.Errorf("Failed to get all users, some might be banned. %s", err)
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
		oldConfig := u.GetUserConfig(uData.ID)

		if !oldConfig.isEditor(userID) {
			return api.NewApiError(http.StatusForbidden, fmt.Errorf("user is not editor"))
		}
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
