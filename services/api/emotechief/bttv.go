package emotechief

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type bttvDashboardResponse struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Displayname   string   `json:"displayName"`
	Providerid    string   `json:"providerId"`
	Bots          []string `json:"bots"`
	Channelemotes []struct {
		ID             string    `json:"id"`
		Code           string    `json:"code"`
		Imagetype      string    `json:"imageType"`
		Userid         string    `json:"userId"`
		Createdat      time.Time `json:"createdAt"`
		Updatedat      time.Time `json:"updatedAt"`
		Global         bool      `json:"global"`
		Live           bool      `json:"live"`
		Sharing        bool      `json:"sharing"`
		Approvalstatus string    `json:"approvalStatus"`
	} `json:"channelEmotes"`
	Sharedemotes []struct {
		ID             string    `json:"id"`
		Code           string    `json:"code"`
		Imagetype      string    `json:"imageType"`
		Createdat      time.Time `json:"createdAt"`
		Updatedat      time.Time `json:"updatedAt"`
		Global         bool      `json:"global"`
		Live           bool      `json:"live"`
		Sharing        bool      `json:"sharing"`
		Approvalstatus string    `json:"approvalStatus"`
		User           struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Displayname string `json:"displayName"`
			Providerid  string `json:"providerId"`
		} `json:"user"`
	} `json:"sharedEmotes"`
}

type dashboardsResponse []dashboardCfg

type dashboardCfg struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Displayname string `json:"displayName"`
	Providerid  string `json:"providerId"`
	Avatar      string `json:"avatar"`
	Limits      struct {
		Channelemotes  int `json:"channelEmotes"`
		Sharedemotes   int `json:"sharedEmotes"`
		Personalemotes int `json:"personalEmotes"`
	} `json:"limits"`
}

type accountResponse struct {
	ID                    string      `json:"id"`
	Flags                 int         `json:"flags"`
	Name                  string      `json:"name"`
	Displayname           string      `json:"displayName"`
	Avatar                string      `json:"avatar"`
	Providerid            string      `json:"providerId"`
	Bots                  []string    `json:"bots"`
	Subscriptionid        interface{} `json:"subscriptionId"`
	Subscriptioncreatedat interface{} `json:"subscriptionCreatedAt"`
	Glow                  bool        `json:"glow"`
	Plan                  string      `json:"plan"`
	Limits                struct {
		Channelemotes  int `json:"channelEmotes"`
		Sharedemotes   int `json:"sharedEmotes"`
		Personalemotes int `json:"personalEmotes"`
	} `json:"limits"`
	Discord struct {
		Userid  interface{} `json:"userId"`
		Guildid interface{} `json:"guildId"`
	} `json:"discord"`
	Email string `json:"email"`
}

func (e *EmoteChief) SetEmote(channelUserID, emoteId string) error {
	// first figure out the bttvUserId for the channel, might cache this later on
	resp, err := http.Get("https://api.betterttv.net/3/cached/users/twitch/" + channelUserID)
	if err != nil {
		return err
	}

	var userResp struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		return err
	}
	bttvUserId := userResp.ID

	// figure out the limit for the channel, might also chache this later on with expiry
	req, err := http.NewRequest("GET", "https://api.betterttv.net/3/account/dashboards", nil)
	req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
	if err != nil {
		log.Error(err)
	}

	resp, err = e.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	var dashboards dashboardsResponse
	err = json.NewDecoder(resp.Body).Decode(&dashboards)
	if err != nil {
		return err
	}

	var dbCfg dashboardCfg
	for _, db := range dashboards {
		if db.ID == bttvUserId {
			dbCfg = db
		}
	}

	sharedEmotesLimit := 0
	if dbCfg.ID == "" {
		// the own account is fetched differently
		req, err := http.NewRequest("GET", "https://api.betterttv.net/3/account", nil)
		req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
		if err != nil {
			log.Error(err)
		}

		resp, err = e.httpClient.Do(req)
		if err != nil {
			log.Error(err)
			return err
		}

		var account accountResponse
		err = json.NewDecoder(resp.Body).Decode(&account)
		if err != nil {
			return err
		}

		if account.ID == bttvUserId {
			sharedEmotesLimit = account.Limits.Sharedemotes
		} else {
			return errors.New("Dashboard not found in account, no permission to moderate")
		}
	} else {
		sharedEmotesLimit = dbCfg.Limits.Sharedemotes
	}

	// figure currently added emotes
	resp, err = http.Get("https://api.betterttv.net/3/users/" + bttvUserId + "?limited=false&personal=false")
	if err != nil {
		return err
	}

	var dashboard bttvDashboardResponse
	err = json.NewDecoder(resp.Body).Decode(&dashboard)
	if err != nil {
		return err
	}

	// figure out the current emote
	currentEmoteId := e.store.Client.HGet("bttv_emote", channelUserID).Val()
	if len(dashboard.Sharedemotes) >= sharedEmotesLimit || currentEmoteId != "" {
		if currentEmoteId == "" || len(dashboard.Sharedemotes) > 0 {
			currentEmoteId = dashboard.Sharedemotes[rand.Intn(len(dashboard.Sharedemotes))].ID
		}

		// Delete the current emote
		req, err = http.NewRequest("DELETE", "https://api.betterttv.net/3/emotes/"+currentEmoteId+"/shared/"+bttvUserId, nil)
		req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
		if err != nil {
			log.Error(err)
		}

		resp, err = e.httpClient.Do(req)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Infof("Deleted: %s %s %d", bttvUserId, currentEmoteId, resp.StatusCode)
	}

	// Add new emote
	req, err = http.NewRequest("PUT", "https://api.betterttv.net/3/emotes/"+emoteId+"/shared/"+bttvUserId, nil)
	req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
	if err != nil {
		log.Error(err)
		return err
	}

	resp, err = e.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Infof("Added: %s %s %d", bttvUserId, emoteId, resp.StatusCode)

	return nil
}
