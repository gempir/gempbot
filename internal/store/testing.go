package store

import "github.com/gempir/gempbot/internal/dto"

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
