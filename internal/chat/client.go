package chat

import (
	"time"

	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/go-twitch-irc/v3"
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

func (c *ChatClient) Reply(channel string, parentMsgId string, message string) {
	c.ircClient.Reply(channel, parentMsgId, message)
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
	select {
	case <-c.connected:
		log.Info("Twitch irc connection established")
		return
	case <-time.After(5 * time.Second):
		log.Info("Bot connection timed out")
		return
	}
}

func (c *ChatClient) Connect(onConnect func()) {
	c.ircClient.OnConnect(func() {
		log.Info("connected to Twitch IRC")
		go func() {
			c.connected <- true
		}()
		onConnect()
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
