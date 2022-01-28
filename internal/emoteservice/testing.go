package emoteservice

type MockApiClient struct {
}

func NewMockApiClient() *MockApiClient {
	return &MockApiClient{}
}

func (c *MockApiClient) GetEmote(emoteID string) (Emote, error) {
	return Emote{}, nil
}

func (c *MockApiClient) RemoveEmote(channelID string, emoteID string) error {
	return nil
}

func (c *MockApiClient) AddEmote(channelID, emoteID string) error {
	return nil
}

func (c *MockApiClient) GetUser(channelID string) (User, error) {
	return User{EmoteSlots: 100}, nil
}
