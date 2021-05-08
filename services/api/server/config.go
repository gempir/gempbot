package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type UserConfig struct {
	Redemptions Redemptions
	Editors     []string
	Protected   Protected
}

type Protected struct {
	EditorFor []string
}

func createDefaultUserConfig() UserConfig {
	return UserConfig{
		Redemptions: Redemptions{
			Bttv: Redemption{Title: "Bttv emote", Active: false},
		},
		Editors: []string{},
		Protected: Protected{
			EditorFor: []string{},
		},
	}
}

func (s *Server) handleUserConfig(w http.ResponseWriter, r *http.Request) {
	ok, auth, _ := s.authenticate(r)
	if !ok {
		http.Error(w, "bad authentication", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		userConfig, err := s.getUserConfig(auth.Data.UserID)
		if err != nil {
			http.Error(w, "can't recover config"+err.Error(), http.StatusBadRequest)
			return
		}

		managing := r.URL.Query().Get("managing")
		if managing != "" {
			userData, err := s.helixClient.GetUsersByUsernames([]string{managing})
			if err != nil || len(userData) == 0 {
				http.Error(w, "can't resolve managing in config "+err.Error(), http.StatusBadRequest)
			}

			isEditor := false
			for _, editor := range userConfig.Protected.EditorFor {
				if editor == userData[managing].ID {
					isEditor = true
				}
			}

			if !isEditor {
				http.Error(w, "User is not editor", http.StatusBadRequest)
			}

			managingUserConfig, err := s.getUserConfig(userData[managing].ID)
			if err != nil {
				http.Error(w, "can't recover config"+err.Error(), http.StatusBadRequest)
				return
			}
			newEditorForNames := []string{}

			userData, err = s.helixClient.GetUsersByUserIds(userConfig.Protected.EditorFor)
			if err != nil {
				http.Error(w, "can't resolve editorFor in config "+err.Error(), http.StatusBadRequest)
			}
			for _, user := range userData {
				newEditorForNames = append(newEditorForNames, user.Login)
			}

			managingUserConfig.Protected.EditorFor = newEditorForNames
			writeJSON(w, managingUserConfig, http.StatusOK)
		} else {
			newEditorNames := []string{}

			userData, err := s.helixClient.GetUsersByUserIds(userConfig.Editors)
			if err != nil {
				http.Error(w, "can't resolve editors in config "+err.Error(), http.StatusBadRequest)
			}
			for _, user := range userData {
				newEditorNames = append(newEditorNames, user.Login)
			}

			userConfig.Editors = newEditorNames

			newEditorForNames := []string{}

			userData, err = s.helixClient.GetUsersByUserIds(userConfig.Protected.EditorFor)
			if err != nil {
				http.Error(w, "can't resolve editorFor in config "+err.Error(), http.StatusBadRequest)
			}
			for _, user := range userData {
				newEditorForNames = append(newEditorForNames, user.Login)
			}

			userConfig.Protected.EditorFor = newEditorForNames
			writeJSON(w, userConfig, http.StatusOK)
		}

	} else if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Errorf("Failed reading update body: %s", err)
			http.Error(w, "Failure saving body: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.processConfig(auth.Data.UserID, body)
		if err != nil {
			log.Errorf("failed processing config: %s", err)
			http.Error(w, "failed processing config: "+err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, "", http.StatusOK)
	} else if r.Method == http.MethodDelete {
		_, err := s.store.Client.HDel("userConfig", auth.Data.UserID).Result()
		if err != nil {
			log.Error(err)
			http.Error(w, "Failed deleting: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = s.unsubscribeChannelPoints(auth.Data.UserID)
		if err != nil {
			log.Error(err)
			http.Error(w, "Failed to unsubscribe"+err.Error(), http.StatusInternalServerError)
			return
		}

		writeJSON(w, "", http.StatusOK)
	}

}

func (s *Server) getUserConfig(userID string) (UserConfig, error) {
	val, err := s.store.Client.HGet("userConfig", userID).Result()
	if err != nil || val == "" {
		return createDefaultUserConfig(), nil
	}

	var userConfig UserConfig
	if err := json.Unmarshal([]byte(val), &userConfig); err != nil {
		return UserConfig{}, errors.New("can't find config")
	}

	return userConfig, nil
}

func (s *Server) processConfig(userID string, body []byte) error {
	isNew := false

	val, err := s.store.Client.HGet("userConfig", userID).Result()
	if err == redis.Nil {
		isNew = true
	} else if err != nil {
		return err
	}

	var oldConfig UserConfig
	if !isNew {
		if err := json.Unmarshal([]byte(val), &oldConfig); err != nil {
			return err
		}
	}

	var newConfig UserConfig
	if err := json.Unmarshal(body, &newConfig); err != nil {
		return err
	}

	protected := oldConfig.Protected
	if protected.EditorFor != nil {
		protected.EditorFor = []string{}
	}

	newEditorIds := []string{}

	userData, err := s.helixClient.GetUsersByUsernames(newConfig.Editors)
	if err != nil {
		return err
	}
	if len(newConfig.Editors) != len(userData) {
		return errors.New("Failed to find all editors")
	}

	for _, user := range userData {
		if user.ID == userID {
			return errors.New("You can't be your own editor")
		}

		newEditorIds = append(newEditorIds, user.ID)
	}

	configToSave := UserConfig{
		Redemptions: newConfig.Redemptions,
		Editors:     newEditorIds,
		Protected:   protected,
	}

	js, err := json.Marshal(configToSave)
	if err != nil {
		return err
	}

	_, err = s.store.Client.HSet("userConfig", userID, js).Result()
	if err != nil {
		return err
	}

	for _, user := range userData {
		err := s.addEditorFor(user.ID, userID)
		if err != nil {
			return err
		}
	}

	if isNew {
		log.Info("Created new config for: ", userID)
		s.subscribeChannelPoints(userID)
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data interface{}, code int) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err = w.Write(js)
	if err != nil {
		log.Errorf("Faile to writeJSON: %s", err)
	}
}

func (s *Server) addEditorFor(editorID, userID string) error {
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
