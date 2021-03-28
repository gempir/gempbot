package collector

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/services/bot/helix"
	"github.com/gempir/spamchamp/services/bot/store"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	startTime    time.Time
	cfg          *config.Config
	helixClient  *helix.Client
	twitchClient *twitch.Client
	store        *store.Store
	channels     map[string]helix.UserData
	joined       map[string]bool
}

// NewBot create new bot instance
func NewBot(cfg *config.Config, helixClient *helix.Client, store *store.Store) *Bot {
	channels, err := helixClient.GetUsersByUserIds(cfg.Channels)
	if err != nil {
		log.Fatalf("[collector] failed to load configured channels %s", err.Error())
	}

	return &Bot{
		cfg:         cfg,
		helixClient: helixClient,
		store:       store,
		channels:    channels,
		joined:      map[string]bool{},
	}
}

// Connect startup the logger and bot
func (b *Bot) Connect() {
	b.startTime = time.Now()
	b.twitchClient = twitch.NewClient(b.cfg.Username, "oauth:"+b.cfg.OAuth)
	b.twitchClient.IrcAddress = "127.0.0.1:3333"
	b.twitchClient.TLS = false
	b.twitchClient.SetupCmd = "LOGIN spamchampbot"

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("[collector] joining as anonymous")
	} else {
		log.Info("[collector] joining as user " + b.cfg.Username)
	}
	b.initialJoins()

	b.twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		b.handlePrivateMessage(message)
	})

	go func() {
		ticker := time.NewTicker(15 * time.Minute)

		for range ticker.C {
			b.LoadTopChannelsAndJoin()
		}
	}()

	log.Fatal(b.twitchClient.Connect())
}

func (b *Bot) initialJoins() {
	for _, channel := range b.channels {
		log.Info("[collector] joining " + channel.Login)
		b.twitchClient.Join(channel.Login)
	}
}

func (b *Bot) joinStoreChannels() {
	go func() {
		allChannelIds := b.store.GetAllChannels()
		channels, err := b.helixClient.GetUsersByUserIds(allChannelIds)
		if err != nil {
			log.Error(err)
			return
		}

		go func() {
			for _, channelID := range allChannelIds {
				if _, ok := channels[channelID]; !ok {
					b.store.RemoveChannel(channelID)
				}
			}
		}()

		for _, userData := range channels {
			if _, ok := b.joined[userData.Login]; !ok {
				log.Info(userData.Login)
				b.joined[userData.Login] = true
				b.twitchClient.Join(userData.Login)
				log.Debugf("[collector] joined %s", userData.DisplayName)
			}
		}
	}()
}

func (b *Bot) SaveAndJoinChannel(channelName string) {
	users, err := b.helixClient.GetUsersByUsernames([]string{channelName})
	if err != nil {
		log.Error(err)
		return
	}

	userids := []string{}
	for _, channel := range users {
		userids = append(userids, channel.ID)
	}

	b.store.AddChannels(userids...)
	b.joinStoreChannels()
}

func (b *Bot) LoadTopChannelsAndJoin() {
	log.Info("[collector] fetching top channels and joining")
	b.store.AddChannels(b.helixClient.GetTopChannels()...)

	b.joinStoreChannels()
}
