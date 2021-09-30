package chat

import (
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/log"
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

func (c *ChatClient) JoinBot(channel string) {
	c.Say(c.cfg.Username, "JOIN "+c.cfg.Environment+" "+channel)
}

func (c *ChatClient) PartBot(channel string) {
	c.Say(c.cfg.Username, "PART "+c.cfg.Environment+" "+channel)
}

func (c *ChatClient) Connect() {
	go func() {
		err := c.ircClient.Connect()
		if err != nil {
			log.Error()
			c.connected <- false
		}
	}()

	go func() {
		c.ircClient.OnConnect(func() {
			c.connected <- true
		})
	}()
}
