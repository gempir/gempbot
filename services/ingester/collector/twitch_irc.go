package collector

import (
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
	exit        chan string
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
		exit:        make(chan string),
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

	<-b.exit
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
				b.joined.mutex.Lock()
				b.joined.m[user.Login] = true
				b.joined.mutex.Unlock()
				b.scaler.Join(user.Login)
				log.Infof("joined %s", user.Login)
			}
		}
	}()
}
