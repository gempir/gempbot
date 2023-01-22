package dto

import "github.com/gempir/go-twitch-irc/v4"

type CommandPayload struct {
	Query string
	Name  string
	Msg   twitch.PrivateMessage
}

const (
	CmdNamePrediction = "prediction"
	CmdNameStatus     = "status"
	CmdNameOutcome    = "outcome"
)
