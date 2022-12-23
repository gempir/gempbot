package helixclient

import (
	"github.com/gempir/gempbot/internal/store"
	"github.com/nicklaw5/helix/v2"
)

type MockHelixClient struct{}

func NewMockClient() *MockHelixClient {
	return &MockHelixClient{}
}

func (m *MockHelixClient) StartRefreshTokenRoutine() {}

func (m *MockHelixClient) RefreshToken(token store.UserAccessToken) error {
	return nil
}

func (m *MockHelixClient) GetTopChannels() []string {
	return []string{}
}

func (m *MockHelixClient) CreateEventSubSubscription(userID string, webHookUrl string, subType string) (*helix.EventSubSubscriptionsResponse, error) {
	return nil, nil
}

func (m *MockHelixClient) CreateRewardEventSubSubscription(userID, webHookUrl, subType, rewardID string, retry bool) (*helix.EventSubSubscriptionsResponse, error) {
	return nil, nil
}

func (m *MockHelixClient) RemoveEventSubSubscription(id string) (*helix.RemoveEventSubSubscriptionParamsResponse, error) {
	return nil, nil
}

func (m *MockHelixClient) GetEventSubSubscriptions(params *helix.EventSubSubscriptionsParams) (*helix.EventSubSubscriptionsResponse, error) {
	return nil, nil
}

func (m *MockHelixClient) GetAllSubscriptions(eventType string) []helix.EventSubSubscription {
	return []helix.EventSubSubscription{}
}

func (m *MockHelixClient) GetPredictions(params *helix.PredictionsParams) (*helix.PredictionsResponse, error) {
	return nil, nil
}

func (m *MockHelixClient) EndPrediction(params *helix.EndPredictionParams) (*helix.PredictionsResponse, error) {
	return nil, nil
}

func (m *MockHelixClient) CreatePrediction(params *helix.CreatePredictionParams) (*helix.PredictionsResponse, error) {
	return nil, nil
}

func (m *MockHelixClient) CreateOrUpdateReward(userID, userAccessToken string, reward CreateCustomRewardRequest, rewardID string) (*helix.ChannelCustomReward, error) {
	return &helix.ChannelCustomReward{}, nil
}

func (m *MockHelixClient) UpdateRedemptionStatus(broadcasterID, rewardID string, redemptionID string, statusSuccess bool) error {
	return nil
}

func (m *MockHelixClient) DeleteReward(userID string, userAccessToken string, rewardID string) error {
	return nil
}

func (m *MockHelixClient) GetUsersByUserIds(userIDs []string) (map[string]UserData, error) {
	return nil, nil
}

func (m *MockHelixClient) GetUsersByUsernames(usernames []string) (map[string]UserData, error) {
	return nil, nil
}

func (m *MockHelixClient) GetUserByUsername(username string) (UserData, error) {
	return UserData{
		ID:    "123",
		Login: "testusergempir",
	}, nil
}

func (m *MockHelixClient) GetUserByUserID(userID string) (UserData, error) {
	return UserData{}, nil
}

func (m *MockHelixClient) SetUserAccessToken(token string) {}

func (m *MockHelixClient) ValidateToken(accessToken string) (bool, *helix.ValidateTokenResponse, error) {
	return false, nil, nil
}

func (m *MockHelixClient) RequestUserAccessToken(code string) (*helix.UserAccessTokenResponse, error) {
	return nil, nil
}
