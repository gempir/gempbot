package emoteservice

type sevenTvEmote struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner struct {
		ID          string `json:"id"`
		TwitchID    string `json:"twitch_id"`
		Login       string `json:"login"`
		DisplayName string `json:"display_name"`
		Role        struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Position int    `json:"position"`
			Color    int    `json:"color"`
			Allowed  int    `json:"allowed"`
			Denied   int    `json:"denied"`
			Default  bool   `json:"default"`
		} `json:"role"`
	} `json:"owner"`
	Visibility       int           `json:"visibility"`
	VisibilitySimple []interface{} `json:"visibility_simple"`
	Mime             string        `json:"mime"`
	Status           int           `json:"status"`
	Tags             []interface{} `json:"tags"`
	Width            []int         `json:"width"`
	Height           []int         `json:"height"`
	Urls             [][]string    `json:"urls"`
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
