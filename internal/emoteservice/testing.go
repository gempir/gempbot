package emoteservice

type TestingClient struct {
}

func (c *TestingClient) GetEmote(emoteID string) (Emote, error) {
	return Emote{}, nil
}

func (c *TestingClient) RemoveEmote(channelID string, emoteID string) error {
	return nil
}

func (c *TestingClient) AddEmote(channelID, emoteID string) error {
	return nil
}

func (c *TestingClient) GetUser(channelID string) (User, error) {
	return User{EmoteSlots: 100}, nil
}
