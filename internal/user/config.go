package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gempir/bot/pkg/api"
	"github.com/gempir/bot/pkg/helix"
	"github.com/gempir/bot/pkg/log"
	"github.com/gempir/bot/pkg/slice"
	"github.com/gempir/bot/pkg/store"
)

type UserAdmin struct {
	db          *store.Database
	helixClient *helix.Client
}

func NewUserAdmin(db *store.Database, helixClient *helix.Client) *UserAdmin {
	return &UserAdmin{db, helixClient}
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

func (u *UserAdmin) convertUserConfig(uCfg UserConfig, toNames bool) (UserConfig, error) {
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
		return UserConfig{}, errors.New("Invalid username(s)")
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

	return uCfg, err
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

func (u *UserAdmin) checkIsEditor(editorUserID string, ownerUserID string) api.Error {
	if editorUserID == ownerUserID {
		return nil
	}

	userConfig := u.GetUserConfig(ownerUserID)

	if userConfig.isEditor(editorUserID) {
		return nil
	}

	return api.NewApiError(http.StatusForbidden, fmt.Errorf("user is not editor"))
}
