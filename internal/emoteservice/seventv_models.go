package emoteservice

type sevenTvEmote struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Flags     int      `json:"flags"`
	Tags      []string `json:"tags"`
	Lifecycle int      `json:"lifecycle"`
	State     []string `json:"state"`
	Listed    bool     `json:"listed"`
	Animated  bool     `json:"animated"`
	Owner     struct {
		ID          string `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
		Style       struct {
			Color int `json:"color"`
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
			FrameCount int    `json:"frame_count"`
			Size       int    `json:"size"`
			Format     string `json:"format"`
		} `json:"files"`
	} `json:"host"`
	Versions []struct {
		ID          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Lifecycle   int      `json:"lifecycle"`
		State       []string `json:"state"`
		Listed      bool     `json:"listed"`
		Animated    bool     `json:"animated"`
		CreatedAt   int64    `json:"createdAt"`
	} `json:"versions"`
}

type gqlQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type SevenTvUserResponse struct {
	Data struct {
		User struct {
			ID     string `json:"id"`
			Emotes []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Status     int    `json:"status"`
				Visibility int    `json:"visibility"`
				Width      []int  `json:"width"`
				Height     []int  `json:"height"`
			} `json:"emotes"`
			EmoteSlots int `json:"emote_slots"`
		} `json:"user"`
	} `json:"data"`
}

const (
	EmoteVisibilityPrivate int32 = 1 << iota
	EmoteVisibilityGlobal
	EmoteVisibilityUnlisted
	EmoteVisibilityOverrideBTTV
	EmoteVisibilityOverrideFFZ
	EmoteVisibilityOverrideTwitchGlobal
	EmoteVisibilityOverrideTwitchSubscriber
	EmoteVisibilityZeroWidth
	EmoteVisibilityPermanentlyUnlisted

	EmoteVisibilityAll int32 = (1 << iota) - 1
)
