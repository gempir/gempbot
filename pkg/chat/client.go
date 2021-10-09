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
	c.ircClient.Say(channel, message)
}

func (c *ChatClient) JoinBot(channel string) {
	c.Say(c.cfg.Username, "JOIN "+c.cfg.Environment+" "+channel)
}

func (c *ChatClient) PartBot(channel string) {
	c.Say(c.cfg.Username, "PART "+c.cfg.Environment+" "+channel)
}

func (c *ChatClient) Join(channel string) {
	log.Infof("JOIN %s", channel)
	c.ircClient.Join(channel)
}

func (c *ChatClient) Part(channel string) {
	log.Infof("PART %s", channel)
	c.ircClient.Depart(channel)
}

func (c *ChatClient) WaitForConnect() {
	<-c.connected
}

func (c *ChatClient) Connect() {
	c.ircClient.OnConnect(func() {
		log.Info("connected to Twitch IRC")
		go func() {
			c.connected <- true
		}()
	})

	count := 0
	c.ircClient.OnRoomStateMessage(func(roomStateMessage twitch.RoomStateMessage) {
		count++
		log.Infof("%d #%s roomstate %v", count, roomStateMessage.Channel, roomStateMessage.State)
	})

	err := c.ircClient.Connect()
	if err != nil {
		log.Error(err)
	}
}

func (c *ChatClient) SetOnPrivateMessage(f func(twitch.PrivateMessage)) {
	c.ircClient.OnPrivateMessage(f)
}
