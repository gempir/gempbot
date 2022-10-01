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

type ConnectionResponse struct {
	ID            string `json:"id"`
	Platform      string `json:"platform"`
	Username      string `json:"username"`
	DisplayName   string `json:"display_name"`
	LinkedAt      int64  `json:"linked_at"`
	EmoteCapacity int    `json:"emote_capacity"`
	EmoteSet      struct {
		ID         string        `json:"id"`
		Name       string        `json:"name"`
		Tags       []interface{} `json:"tags"`
		Immutable  bool          `json:"immutable"`
		Privileged bool          `json:"privileged"`
		Emotes     []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			Flags     int    `json:"flags"`
			Timestamp int64  `json:"timestamp"`
			ActorID   string `json:"actor_id"`
			Data      struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				Flags     int    `json:"flags"`
				Lifecycle int    `json:"lifecycle"`
				Listed    bool   `json:"listed"`
				Animated  bool   `json:"animated"`
				Owner     struct {
					ID          string `json:"id"`
					Username    string `json:"username"`
					DisplayName string `json:"display_name"`
					Style       struct {
						Color int         `json:"color"`
						Paint interface{} `json:"paint"`
					} `json:"style"`
					Roles []string `json:"roles"`
				} `json:"owner"`
				Host struct {
					URL   string `json:"url"`
					Files []struct {
						Name       string `json:"name"`
						StaticName string `json:"static_name"`
						Width      int    `json:"width"`
						Height     int    `json:"height"`
						Size       int    `json:"size"`
						Format     string `json:"format"`
					} `json:"files"`
				} `json:"host"`
			} `json:"data"`
		} `json:"emotes"`
		Capacity int `json:"capacity"`
		Owner    struct {
			ID          string `json:"id"`
			Username    string `json:"username"`
			DisplayName string `json:"display_name"`
			Style       struct {
				Color int         `json:"color"`
				Paint interface{} `json:"paint"`
			} `json:"style"`
			Roles []string `json:"roles"`
		} `json:"owner"`
	} `json:"emote_set"`
	User struct {
		ID                string `json:"id"`
		Username          string `json:"username"`
		ProfilePictureURL string `json:"profile_picture_url"`
		DisplayName       string `json:"display_name"`
		Style             struct {
			Color int         `json:"color"`
			Paint interface{} `json:"paint"`
		} `json:"style"`
		Biography string `json:"biography"`
		Editors   []struct {
			ID          string `json:"id"`
			Permissions int    `json:"permissions"`
			Visible     bool   `json:"visible"`
			AddedAt     int64  `json:"added_at"`
		} `json:"editors"`
		Roles       []string `json:"roles"`
		Connections []struct {
			ID            string `json:"id"`
			Platform      string `json:"platform"`
			Username      string `json:"username"`
			DisplayName   string `json:"display_name"`
			LinkedAt      int64  `json:"linked_at"`
			EmoteCapacity int    `json:"emote_capacity"`
		} `json:"connections"`
	} `json:"user"`
}
