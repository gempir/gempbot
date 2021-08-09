package emotechief

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gempir/bitraft/pkg/log"
	"github.com/gempir/bitraft/pkg/store"
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

type bttvEmoteResponse struct {
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

func (e *EmoteChief) SetBttvEmote(channelUserID, emoteId, channel string, slots int) (addedEmote *bttvEmoteResponse, removedEmote *bttvEmoteResponse, err error) {
	addedEmote, err = getBttvEmote(emoteId)
	if err != nil {
		return
	}

	if !addedEmote.Sharing {
		err = errors.New("Emote is not shared")
		return
	}

	// first figure out the bttvUserId for the channel, might cache this later on
	var resp *http.Response
	resp, err = http.Get("https://api.betterttv.net/3/cached/users/twitch/" + channelUserID)
	if err != nil {
		return
	}

	var userResp struct {
		ID string `json:"id"`
	}
	err = json.NewDecoder(resp.Body).Decode(&userResp)
	if err != nil {
		return
	}
	bttvUserId := userResp.ID

	// figure out the limit for the channel, might also chache this later on with expiry
	var req *http.Request
	req, err = http.NewRequest("GET", "https://api.betterttv.net/3/account/dashboards", nil)
	req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err = e.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return
	}

	var dashboards dashboardsResponse
	err = json.NewDecoder(resp.Body).Decode(&dashboards)
	if err != nil {
		return
	}

	var dbCfg dashboardCfg
	for _, db := range dashboards {
		if db.ID == bttvUserId {
			dbCfg = db
		}
	}
	if dbCfg.ID == "" {
		err = errors.New("No permission to moderate, add gempbot as BetterTTV editor")
		return
	}
	sharedEmotesLimit := dbCfg.Limits.Sharedemotes

	// figure currently added emotes
	resp, err = http.Get("https://api.betterttv.net/3/users/" + bttvUserId + "?limited=false&personal=false")
	if err != nil {
		return
	}

	var dashboard bttvDashboardResponse
	err = json.NewDecoder(resp.Body).Decode(&dashboard)
	if err != nil {
		return
	}

	for _, emote := range dashboard.Sharedemotes {
		if emote.ID == emoteId {
			err = errors.New("Emote already added")
			return
		}
		if emote.Code == addedEmote.Code {
			err = errors.New("Emote code already added")
			return
		}
	}

	for _, emote := range dashboard.Channelemotes {
		if emote.ID == emoteId {
			err = errors.New("Emote already a channelEmote")
			return
		}
		if emote.Code == addedEmote.Code {
			err = errors.New("Emote code already a channelEmote")
			return
		}
	}
	log.Debugf("Current shared emotes: %d/%d", len(dashboard.Sharedemotes), sharedEmotesLimit)

	var removalTargetEmoteId string

	emotesAdded := e.db.GetEmoteAdded(channelUserID, slots)
	log.Infof("Total Previous emotes %d in %s", len(emotesAdded), channelUserID)

	confirmedEmotesAdded := []store.EmoteAdd{}
	for _, emote := range emotesAdded {
		for _, sharedEmote := range dashboard.Sharedemotes {
			if emote.EmoteID == sharedEmote.ID {
				confirmedEmotesAdded = append(confirmedEmotesAdded, emote)
			}
		}
	}

	if len(confirmedEmotesAdded) == slots {
		removalTargetEmoteId = confirmedEmotesAdded[len(confirmedEmotesAdded)-1].EmoteID
		log.Infof("Found removal target %s in %s", removalTargetEmoteId, channelUserID)
	} else if len(dashboard.Sharedemotes) >= sharedEmotesLimit {
		log.Infof("Didn't find previous emote history of %d emotes and limit reached, choosing random in %s", slots, channelUserID)
		removalTargetEmoteId = dashboard.Sharedemotes[rand.Intn(len(dashboard.Sharedemotes))].ID
	}

	// do we need to remove the emote?
	if removalTargetEmoteId != "" {
		// Delete the current emote
		req, err = http.NewRequest("DELETE", "https://api.betterttv.net/3/emotes/"+removalTargetEmoteId+"/shared/"+bttvUserId, nil)
		req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
		if err != nil {
			log.Error(err)
			return
		}

		resp, err = e.httpClient.Do(req)
		if err != nil {
			log.Error(err)
			return
		}
		log.Infof("[%d] Deleted channelId: %s emoteId: %s", resp.StatusCode, channelUserID, removalTargetEmoteId)
	}

	// Add new emote
	req, err = http.NewRequest("PUT", "https://api.betterttv.net/3/emotes/"+emoteId+"/shared/"+bttvUserId, nil)
	req.Header.Set("authorization", "Bearer "+e.cfg.BttvToken)
	if err != nil {
		log.Error(err)
		return
	}

	resp, err = e.httpClient.Do(req)
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("[%d] Added channelId: %s emoteId: %s", resp.StatusCode, channelUserID, emoteId)

	if resp.StatusCode < http.StatusBadRequest {
		e.db.CreateEmoteAdd(channelUserID, emoteId)
	}

	if removalTargetEmoteId != "" {
		removedEmote, err = getBttvEmote(removalTargetEmoteId)
		if err != nil {
			return
		}
	}

	return
}

func getBttvEmote(emoteID string) (*bttvEmoteResponse, error) {
	if emoteID == "" {
		return nil, nil
	}

	response, err := http.Get("https://api.betterttv.net/3/emotes/" + emoteID)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if response.StatusCode <= 100 || response.StatusCode >= 400 {
		return nil, fmt.Errorf("Bad bttv response: %d", response.StatusCode)
	}

	var emoteResponse bttvEmoteResponse
	err = json.NewDecoder(response.Body).Decode(&emoteResponse)
	if err != nil {
		return nil, err
	}

	return &emoteResponse, nil
}
