package collector

import (
	"strings"
	"sync"
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
	channels    stringUserDataSyncMap
	joined      stringBoolSyncMap
	active      stringBoolSyncMap
}

type stringUserDataSyncMap struct {
	m     map[string]helix.UserData
	mutex *sync.Mutex
}

type stringBoolSyncMap struct {
	m     map[string]bool
	mutex *sync.Mutex
}

// NewBot create new bot instance
func NewBot(cfg *config.Config, helixClient *helix.Client, store *store.Store) *Bot {
	channels, err := helixClient.GetUsersByUserIds(cfg.Channels)
	if err != nil {
		log.Fatalf("[collector] failed to load configured channels %s", err.Error())
	}

	channelsMap := stringUserDataSyncMap{m: map[string]helix.UserData{}, mutex: &sync.Mutex{}}
	for key, data := range channels {
		channelsMap.mutex.Lock()
		channelsMap.m[key] = data
		channelsMap.mutex.Unlock()
	}

	return &Bot{
		cfg:         cfg,
		helixClient: helixClient,
		store:       store,
		channels:    channelsMap,
		joined:      stringBoolSyncMap{m: map[string]bool{}, mutex: &sync.Mutex{}},
		active:      stringBoolSyncMap{m: map[string]bool{}, mutex: &sync.Mutex{}},
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
		b.store.PublishJoinedChannels(len(b.joined.m))
		b.store.PublishActiveChannels(len(b.active.m))
	}
}

func (b *Bot) initialJoins() {
	for _, channel := range b.channels.m {
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
			if _, ok := b.joined.m[userData.Login]; !ok {
				b.joined.mutex.Lock()
				b.joined.m[userData.Login] = true
				b.joined.mutex.Unlock()
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
