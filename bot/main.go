package main

import (
	"github.com/gempir/go-twitch-irc/v2"
)

var broadcast = make(chan socketMessage)

type socketMessage struct {
	Channel string       `json:"channel"`
	Emote   twitch.Emote `json:"emote"`
	Message string       `json:"message"`
}

func main() {
	client := twitch.NewClient("justinfan123123", "oauth:123123123")

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		msg := &socketMessage{
			Channel: message.Channel,
			Message: message.Message,
		}

		broadcast <- *msg

		if len(message.Emotes) == 0 {
			return
		}
	})

	client.Join("gempir")

	go client.Connect()

	startWebsocketServer()
}
