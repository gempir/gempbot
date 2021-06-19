package collector

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/helix"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/gempir/bitraft/services/ingester/scaler"
	"github.com/gempir/go-twitch-irc/v2"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	startTime   time.Time
	cfg         *config.Config
	scaler      *scaler.Scaler
	store       *store.Redis
	db          *store.Database
	helixClient *helix.Client
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

func NewBot(cfg *config.Config, store *store.Redis, db *store.Database, helixClient *helix.Client) *Bot {
	channelsMap := stringUserDataSyncMap{m: map[string]helix.UserData{}, mutex: &sync.Mutex{}}

	return &Bot{
		cfg:         cfg,
		store:       store,
		db:          db,
		helixClient: helixClient,
		channels:    channelsMap,
		joined:      stringBoolSyncMap{m: map[string]bool{}, mutex: &sync.Mutex{}},
		active:      stringBoolSyncMap{m: map[string]bool{}, mutex: &sync.Mutex{}},
	}
}

func (b *Bot) Connect() {
	b.startTime = time.Now()
	b.scaler = scaler.NewScaler(b.cfg, b.handlePrivateMessage)

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("[collector] joining as anonymous")
	} else {
		log.Info("[collector] joining as user " + b.cfg.Username)
	}
	b.joinBotConfigChannels()

	topic := b.store.SubscribeIngesterMessage()
	channel := topic.Channel()
	for msg := range channel {
		var message store.IngesterMessage

		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Error(err)
		}

		switch message.Type {
		case store.IngesterMsgJoin:
			b.Join(message.Argument)
		case store.IngesterMsgPart:
			b.Part(message.Argument)
		}
	}
}

func (b *Bot) handlePrivateMessage(message twitch.PrivateMessage) {
	b.store.PublishPrivateMessage(message.Raw)
}

func (b *Bot) joinBotConfigChannels() {
	go func() {
		botConfigs := b.db.GetAllJoinBotConfigs()
		userIDs := []string{}
		for _, botConfig := range botConfigs {
			userIDs = append(userIDs, botConfig.OwnerTwitchID)
		}

		users, err := b.helixClient.GetUsersByUserIds(userIDs)
		if err != nil {
			log.Error(err)
		}

		for _, user := range users {
			if _, ok := b.joined.m[user.Login]; !ok {
				b.Join(user.Login)
			}
		}
	}()
}

func (b *Bot) Join(channel string) {
	b.joined.mutex.Lock()
	b.joined.m[channel] = true
	b.joined.mutex.Unlock()
	b.scaler.Join(channel)
	log.Infof("joined %s", channel)
}

func (b *Bot) Part(channel string) {
	b.joined.mutex.Lock()
	delete(b.joined.m, channel)
	b.joined.mutex.Unlock()
	b.scaler.Part(channel)
	log.Infof("part %s", channel)
}
