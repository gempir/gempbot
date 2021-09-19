package chat

import (
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/go-twitch-irc/v2"
)

type ChatClient struct {
	ircClient *twitch.Client
	cfg       *config.Config
}

func NewClient(cfg *config.Config) *ChatClient {
	return &ChatClient{
		cfg:       cfg,
		ircClient: twitch.NewClient(cfg.Username, cfg.OAuth),
	}
}

func (c *ChatClient) Say(channel string, message string) {
	c.ircClient.Say(channel, message)
}

func (c *ChatClient) Connect(done chan bool) {
	go func() {
		c.ircClient.Connect()
	}()

	go func() {
		c.ircClient.OnConnect(func() {
			done <- true
		})
	}()
}
