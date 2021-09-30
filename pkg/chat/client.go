package chat

import (
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/go-twitch-irc/v2"
)

type ChatClient struct {
	ircClient *twitch.Client
	cfg       *config.Config
	connected chan bool
}

func NewClient(cfg *config.Config) *ChatClient {
	return &ChatClient{
		cfg:       cfg,
		connected: make(chan bool),
		ircClient: twitch.NewClient(cfg.Username, cfg.OAuth),
	}
}

func (c *ChatClient) Say(channel string, message string) {
	<-c.connected
	c.ircClient.Say(channel, message)
}

func (c *ChatClient) Connect() {
	go func() {
		c.ircClient.Connect()
	}()

	go func() {
		c.ircClient.OnConnect(func() {
			c.connected <- true
		})
	}()
}
