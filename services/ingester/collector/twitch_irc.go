package collector

import (
	"strings"
	"time"

	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/pkg/helix"
	"github.com/gempir/spamchamp/pkg/store"
	"github.com/gempir/spamchamp/services/ingester/scaler"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	startTime   time.Time
	cfg         *config.Config
	helixClient *helix.Client
	scaler      *scaler.Scaler
	store       *store.Store
	channels    map[string]helix.UserData
	joined      map[string]bool
	active      map[string]bool
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
		active:      map[string]bool{},
	}
}

// Connect startup the logger and bot
func (b *Bot) Connect() {
	b.startTime = time.Now()
	b.scaler = scaler.NewScaler(b.cfg, b.handlePrivateMessage)

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("[collector] joining as anonymous")
	} else {
		log.Info("[collector] joining as user " + b.cfg.Username)
	}
	b.initialJoins()

	go func() {
		time.Sleep(time.Second)
		b.LoadTopChannelsAndJoin()
	}()

	go func() {
		ticker := time.NewTicker(15 * time.Minute)

		for range ticker.C {
			b.LoadTopChannelsAndJoin()
		}
	}()

	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		b.store.PublishJoinedChannels(len(b.joined))
		b.store.PublishActiveChannels(len(b.active))
	}
}

func (b *Bot) initialJoins() {
	for _, channel := range b.channels {
		b.scaler.Join(channel.Login)
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
				b.joined[userData.Login] = true
				b.scaler.Join(userData.Login)
			}
		}
	}()
}

func (b *Bot) SaveAndJoin(channelName string) {
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
