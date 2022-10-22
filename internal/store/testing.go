package store

import (
	"context"

	"github.com/gempir/gempbot/internal/dto"
)

type MockStore struct {
}

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (s *MockStore) IsEmoteBlocked(channelUserID string, emoteID string, rewardType dto.RewardType) bool {
	return false
}

func (s *MockStore) GetEmoteAdded(channelUserID string, rewardType dto.RewardType, slots int) []EmoteAdd {
	return []EmoteAdd{
		{ID: 1, ChannelTwitchID: "channelid", Type: dto.REWARD_SEVENTV, EmoteID: "emoteid"},
	}
}

func (s *MockStore) CreateEmoteAdd(channelUserId string, rewardType dto.RewardType, emoteID string, changeType dto.EmoteChangeType) {
	// do nothing
}

func (s *MockStore) GetUserAccessToken(userID string) (UserAccessToken, error) {
	return UserAccessToken{}, nil
}

func (s *MockStore) GetAppAccessToken() (AppAccessToken, error) {
	return AppAccessToken{
		AccessToken:  "accesstoken",
		RefreshToken: "refreshtoken",
		Scopes:       "scopes",
	}, nil
}

func (s *MockStore) SaveAppAccessToken(ctx context.Context, accessToken string, refreshToken string, scopes string, expiresIn int) error {
	return nil
}

func (s *MockStore) SaveUserAccessToken(ctx context.Context, ownerId string, accessToken string, refreshToken string, scopes string) error {
	return nil
}

func (s *MockStore) GetAllUserAccessToken() []UserAccessToken {
	return []UserAccessToken{
		{OwnerTwitchID: "31231", AccessToken: "accesstoken", RefreshToken: "refreshtoken", Scopes: "scopes"},
	}
}

func (s *MockStore) GetSevenTvToken(ctx context.Context) string {
	return "7tvApiToken"
}

func (s *MockStore) GetBttvToken(ctx context.Context) string {
	return "BttvToken"
}
