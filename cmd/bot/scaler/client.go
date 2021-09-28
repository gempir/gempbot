package scaler

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/go-twitch-irc/v2"
	log "github.com/sirupsen/logrus"
)

type Scaler struct {
	cfg              *config.Config
	clients          []client
	onPrivateMessage func(message twitch.PrivateMessage)
}

type client struct {
	client *twitch.Client
	joined map[string]bool
}

func NewScaler(cfg *config.Config, onPrivateMessage func(message twitch.PrivateMessage)) *Scaler {
	return &Scaler{
		cfg:              cfg,
		clients:          make([]client, 0),
		onPrivateMessage: onPrivateMessage,
	}
}

func (s *Scaler) Join(channelName string) {
	channelName = strings.ToLower(channelName)

	for _, client := range s.clients {
		if len(client.joined) < 50 {
			client.client.Join(channelName)
			client.joined[channelName] = true
			log.Debugf("joining %s", channelName)
			return
		}
	}

	if len(s.clients) > 0 {
		time.Sleep(time.Second * 10)
	}

	log.Infof("Creating new client, total: %d", len(s.clients)+1)

	newClient := twitch.NewClient(s.cfg.Username, s.cfg.OAuth)
	newClient.OnPrivateMessage(s.onPrivateMessage)
	go func() {
		err := newClient.Connect()
		log.Error(err)
	}()

	s.clients = append(s.clients, client{newClient, make(map[string]bool)})

	s.Join(channelName)
}

func (s *Scaler) Part(channel string) {
	for _, client := range s.clients {
		if _, ok := client.joined[channel]; ok {
			client.client.Depart(channel)
		}
	}
}

func (s *Scaler) HandlePrivateMessage(msg twitch.PrivateMessage) {
	sysMessage := msg.Channel == s.cfg.Username && msg.User.Name == s.cfg.Username
	if sysMessage && strings.HasPrefix(msg.Message, "JOIN ") {
		s.Join(strings.TrimPrefix(msg.Message, "JOIN "))
		return
	}
	if sysMessage && strings.HasPrefix(msg.Message, "PART ") {
		s.Part(strings.TrimPrefix(msg.Message, "PART "))
		return
	}

	s.onPrivateMessage(msg)
}

type Message struct {
	Channel string
	Message string
}

func (s *Scaler) Say(msg Message) {
	randomIndex := rand.Intn(len(s.clients))
	s.clients[randomIndex].client.Say(msg.Channel, msg.Message)
}
