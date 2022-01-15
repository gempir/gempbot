package dto

type RewardType string

const (
	REWARD_BTTV    RewardType = "bttv"
	REWARD_SEVENTV RewardType = "seventv"
	REWARD_TIMEOUT RewardType = "timeout"
)

type EmoteChangeType string

const (
	EMOTE_ADD_ADD              EmoteChangeType = "add"
	EMOTE_ADD_REMOVED_PREVIOUS EmoteChangeType = "remove"
	EMOTE_ADD_REMOVED_RANDOM   EmoteChangeType = "removed_random"
	EMOTE_ADD_REMOVED_BLOCKED  EmoteChangeType = "removed_blocked"
)
