package dto

import "github.com/gempir/go-twitch-irc/v2"

type CommandPayload struct {
	Query string
	Name  string
	Msg   twitch.PrivateMessage
}
