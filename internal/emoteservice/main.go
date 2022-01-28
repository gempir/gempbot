package emoteservice

type Emote struct {
	ID   string
	Code string
}

type User struct {
	ID         string
	Emotes     []Emote
	EmoteSlots int
}

type ApiClient interface {
	GetEmote(emoteID string) (Emote, error)
	RemoveEmote(channelID string, emoteID string) error
	AddEmote(channelID, emoteID string) error
	GetUser(channelID string) (User, error)
}
