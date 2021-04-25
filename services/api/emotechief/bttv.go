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

func (e *EmoteChief) SetEmote(userId, emoteId string) error {
	bearerToken := ""
	bttvUserId := "5590fac14d62b7e22aaac72c"
	resp, err := http.Get("https://api.betterttv.net/3/users/" + bttvUserId + "?limited=false&personal=false")
	if err != nil {
		return err
	}

	var dashboard bttvDashboardResponse
	err = json.NewDecoder(resp.Body).Decode(&dashboard)
	if err != nil {
		return err
	}

	log.Info("%v", dashboard)
	// log.Info(dashboard.Sharedemotes)
	currentEmoteId := e.store.Client.HGet("bttv_emote", userId).Val()
	if currentEmoteId == "" || len(dashboard.Sharedemotes) > 0 {
		currentEmoteId = dashboard.Sharedemotes[rand.Intn(len(dashboard.Sharedemotes))].ID
	}

	req, err := http.NewRequest("DELETE", "https://api.betterttv.net/3/emotes/"+bttvUserId+"/shared/"+currentEmoteId, nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
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
	req, err = http.NewRequest("PUT", "https://api.betterttv.net/3/emotes/"+bttvUserId+"/shared/"+emoteId, nil)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
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
