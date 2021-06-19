package scaler

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gempir/bitraft/pkg/config"
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

	log.Infof("Creating new client, previous total: %d", len(s.clients))

	newClient := twitch.NewAnonymousClient()
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

func (s *Scaler) Say(channel, message string) {
	randomIndex := rand.Intn(len(s.clients))
	s.clients[randomIndex].client.Say(channel, message)
}
