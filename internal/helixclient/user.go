package helixclient

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gempir/gempbot/internal/log"
	"github.com/nicklaw5/helix/v2"
)

// UserData exported data from twitch
type UserData struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	ViewCount       int    `json:"view_count"`
	Email           string `json:"email"`
}

func chunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}

// GetUsersByUserIds receive userData for given ids
func (c *Client) GetUsersByUserIds(userIDs []string) (map[string]UserData, error) {
	var filteredUserIDs []string

	for _, id := range userIDs {
		if _, ok := userCacheByID[id]; !ok {
			id = strings.TrimSpace(id)
			if id == "" {
				continue
			}
			filteredUserIDs = append(filteredUserIDs, id)
		}
	}

	if len(filteredUserIDs) > 0 {
		chunks := chunkBy(filteredUserIDs, 50)

		for _, chunk := range chunks {
			resp, err := c.Client.GetUsers(&helix.UsersParams{
				IDs: chunk,
			})
			if err != nil {
				return map[string]UserData{}, err
			}
			if resp.StatusCode == http.StatusUnauthorized {
				return map[string]UserData{}, fmt.Errorf("bad helix response: %v", resp.ErrorMessage)
			}

			log.Infof("[helix] %d GetUsersByUserIds %# v  %# v  %# v ", resp.StatusCode, resp.Error, resp.ErrorMessage, chunk)

			for _, user := range resp.Data.Users {
				data := &UserData{
					ID:              user.ID,
					Login:           user.Login,
					DisplayName:     user.Login,
					Type:            user.Type,
					BroadcasterType: user.BroadcasterType,
					Description:     user.Description,
					ProfileImageURL: user.ProfileImageURL,
					OfflineImageURL: user.OfflineImageURL,
					ViewCount:       user.ViewCount,
					Email:           user.Email,
				}
				cacheMutex.Lock()
				userCacheByID[user.ID] = data
				userCacheByUsername[user.Login] = data
				cacheMutex.Unlock()
			}
		}
	}

	result := make(map[string]UserData)

	for _, id := range userIDs {
		val, ok := userCacheByID[id]
		if !ok {
			log.Infof("Could not find userId, channel might be banned: %s", id)
			continue
		}
		result[id] = *val
	}

	return result, nil
}

// GetUsersByUsernames fetches userdata from helix
func (c *Client) GetUsersByUsernames(usernames []string) (map[string]UserData, error) {
	var filteredUsernames []string

	for _, username := range usernames {
		if _, ok := userCacheByUsername[strings.ToLower(username)]; !ok {
			username = strings.TrimSpace(username)
			if username == "" {
				continue
			}
			filteredUsernames = append(filteredUsernames, strings.ToLower(username))
		}
	}

	if len(filteredUsernames) > 0 {
		chunks := chunkBy(filteredUsernames, 50)

		for _, chunk := range chunks {
			resp, err := c.Client.GetUsers(&helix.UsersParams{
				Logins: chunk,
			})
			if err != nil {
				return map[string]UserData{}, err
			}

			log.Infof("[helix] %d GetUsersByUsernames %# v  %# v  %# v ", resp.StatusCode, resp.Error, resp.ErrorMessage, chunk)
			if resp.StatusCode > http.StatusMultipleChoices {
				return map[string]UserData{}, fmt.Errorf("bad helix response: %v", resp.ErrorMessage)
			}

			for _, user := range resp.Data.Users {
				data := &UserData{
					ID:              user.ID,
					Login:           user.Login,
					DisplayName:     user.Login,
					Type:            user.Type,
					BroadcasterType: user.BroadcasterType,
					Description:     user.Description,
					ProfileImageURL: user.ProfileImageURL,
					OfflineImageURL: user.OfflineImageURL,
					ViewCount:       user.ViewCount,
					Email:           user.Email,
				}
				userCacheByID[user.ID] = data
				userCacheByUsername[user.Login] = data
			}
		}
	}

	result := make(map[string]UserData)

	for _, username := range usernames {
		val, ok := userCacheByUsername[strings.ToLower(username)]
		if !ok {
			log.Infof("Could not find userId, channel might be banned: %s", username)
			continue
		}
		result[strings.ToLower(username)] = *val
	}

	return result, nil
}

func (c *Client) GetUserByUsername(username string) (UserData, error) {
	result, err := c.GetUsersByUsernames([]string{username})
	if err != nil || len(result) != 1 {
		return UserData{}, err
	}

	return result[username], nil
}

func (c *Client) GetUserByUserID(userID string) (UserData, error) {
	result, err := c.GetUsersByUserIds([]string{userID})
	if err != nil || len(result) != 1 {
		return UserData{}, err
	}

	return result[userID], nil
}
