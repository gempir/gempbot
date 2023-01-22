package tmi

import "github.com/gempir/go-twitch-irc/v4"

func IsModerator(user twitch.User) bool {
	val, ok := user.Badges["moderator"]

	if !ok {
		return false
	}

	return val == 1
}

func IsBroadcaster(user twitch.User) bool {
	val, ok := user.Badges["broadcaster"]
	if !ok {
		return false
	}

	return val == 1
}
