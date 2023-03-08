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

func (s *MockStore) AddToQueue(queueItem MediaQueue) error {
	return nil
}

func (s *MockStore) GetQueue(channelTwitchID string) []MediaQueue {
	return []MediaQueue{}
}

func (s *MockStore) GetAllMediaCommandsBotConfig() []BotConfig {
	return []BotConfig{}
}

func (s *MockStore) CreateOrUpdateElection(ctx context.Context, election Election) error {
	return nil
}

func (s *MockStore) GetElection(ctx context.Context, channelTwitchID string) (Election, error) {
	return Election{}, nil
}

func (s *MockStore) DeleteElection(ctx context.Context, channelTwitchID string) error {
	return nil
}

func (s *MockStore) GetAllElections(ctx context.Context) ([]Election, error) {
	return []Election{}, nil
}

func (s *MockStore) SaveReward(reward ChannelPointReward) error {
	return nil
}

func (s *MockStore) CreateOrIncrementNomination(ctx context.Context, nomination Nomination) error {
	return nil
}

func (s *MockStore) GetNominations(ctx context.Context, channelTwitchID string) ([]Nomination, error) {
	return []Nomination{}, nil
}

func (s *MockStore) GetActiveElection(ctx context.Context, channelTwitchID string) (Election, error) {
	return Election{}, nil
}

func (s *MockStore) ClearNominations(ctx context.Context, channelTwitchID string) error {
	return nil
}

func (s *MockStore) DeleteChannelPointRewardById(userID string, rewardID string) {
}

func (s *MockStore) GetChannelPointReward(userID string, rewardType dto.RewardType) (ChannelPointReward, error) {
	return ChannelPointReward{}, nil
}

func (s *MockStore) CreateNominationVote(ctx context.Context, vote NominationVote) error {
	return nil
}

func (s *MockStore) RemoveNominationVote(ctx context.Context, vote NominationVote) error {
	return nil
}

func (s *MockStore) ClearNominationEmote(ctx context.Context, channelTwitchID string, emoteID string) error {
	return nil
}

func (s *MockStore) GetNomination(ctx context.Context, channelTwitchID string, emoteID string) (Nomination, error) {
	return Nomination{}, nil
}

func (s *MockStore) RemoveNomination(ctx context.Context, channelTwitchID string, emoteID string) error {
	return nil
}

func (s *MockStore) CountNominations(ctx context.Context, channelTwitchID string, userID string) (int, error) {
	return 0, nil
}

func (s *MockStore) AddEmoteLogEntry(ctx context.Context, emoteLog EmoteLog) {
}

func (s *MockStore) GetEmoteLogEntries(ctx context.Context, channelTwitchID string, limit int, page int) []EmoteLog {
	return []EmoteLog{}
}

func (s *MockStore) CreateNominationDownvote(ctx context.Context, downvote NominationDownvote) error {
	return nil
}

func (s *MockStore) RemoveNominationDownvote(ctx context.Context, downvote NominationDownvote) error {
	return nil
}

func (s *MockStore) IsAlreadyNominated(ctx context.Context, channelTwitchID string, emoteID string) (bool, error) {
	return false, nil
}

func (s *MockStore) CountNominationDownvotes(ctx context.Context, channelTwitchID string, voteBy string) (int, error) {
	return 0, nil
}

func (s *MockStore) CountNominationVotes(ctx context.Context, channelTwitchID string, voteBy string) (int, error) {
	return 0, nil
}
