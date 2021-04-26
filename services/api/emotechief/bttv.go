package emotechief

import (
	"encoding/json"
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

func (e *EmoteChief) SetEmote(channelUserID, emoteId string) error {
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

	resp, err = http.Get("https://api.betterttv.net/3/users/" + bttvUserId + "?limited=false&personal=false")
	if err != nil {
		return err
	}

	var dashboard bttvDashboardResponse
	err = json.NewDecoder(resp.Body).Decode(&dashboard)
	if err != nil {
		return err
	}

	currentEmoteId := e.store.Client.HGet("bttv_emote", channelUserID).Val()
	if currentEmoteId == "" || len(dashboard.Sharedemotes) > 0 {
		currentEmoteId = dashboard.Sharedemotes[rand.Intn(len(dashboard.Sharedemotes))].ID
	}

	req, err := http.NewRequest("DELETE", "https://api.betterttv.net/3/emotes/"+currentEmoteId+"/shared/"+bttvUserId, nil)
	req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
	req.Header.Set("authority", "api.betterttv.net")
	req.Header.Set("accept", "json")
	if err != nil {
		log.Error(err)
	}

	log.Infof("Deleting %s %s", bttvUserId, currentEmoteId)

	resp, err = e.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("Deletion: %d", resp.StatusCode)

	log.Infof("Adding Emote %s %s", bttvUserId, emoteId)
	req, err = http.NewRequest("PUT", "https://api.betterttv.net/3/emotes/"+emoteId+"/shared/"+bttvUserId, nil)
	req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
	req.Header.Set("authority", "api.betterttv.net")
	req.Header.Set("accept", "json")
	if err != nil {
		log.Error(err)
		return err
	}

	resp, err = e.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Infof("Addition: %d", resp.StatusCode)

	return nil
}
