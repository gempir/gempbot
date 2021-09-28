package collector

import (
	"strings"
	"sync"
	"time"

	"github.com/gempir/gempbot/cmd/bot/commander"
	"github.com/gempir/gempbot/cmd/bot/scaler"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/go-twitch-irc/v2"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	startTime   time.Time
	cfg         *config.Config
	scaler      *scaler.Scaler
	db          *store.Database
	helixClient *helix.Client
	listener    *commander.Listener
	channels    stringUserDataSyncMap
	write       chan scaler.Message
	joined      stringBoolSyncMap
	active      stringBoolSyncMap
	Done        chan bool
}

type stringUserDataSyncMap struct {
	m     map[string]helix.UserData
	mutex *sync.Mutex
}

type stringBoolSyncMap struct {
	m     map[string]bool
	mutex *sync.Mutex
}

func NewBot(cfg *config.Config, db *store.Database, helixClient *helix.Client) *Bot {
	channelsMap := stringUserDataSyncMap{m: map[string]helix.UserData{}, mutex: &sync.Mutex{}}

	write := make(chan scaler.Message)
	handler := commander.NewHandler(cfg, helixClient, db, write)

	listener := commander.NewListener(db, handler, write)
	listener.RegisterDefaultCommands()

	return &Bot{
		Done:        make(chan bool),
		cfg:         cfg,
		db:          db,
		listener:    listener,
		helixClient: helixClient,
		channels:    channelsMap,
		write:       write,
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
	go b.Join(b.cfg.Username)
	go b.joinBotConfigChannels()

	for msg := range b.write {
		b.scaler.Say(scaler.Message{Channel: msg.Channel, Message: msg.Message})
	}
}

func (b *Bot) handlePrivateMessage(msg twitch.PrivateMessage) {
	sysMessage := msg.Channel == b.cfg.Username && msg.User.Name == b.cfg.Username
	if sysMessage {
		log.Infof("sysMessage: %s", msg.Message)
		if strings.HasPrefix(msg.Message, "JOIN ") {
			b.Join(strings.TrimPrefix(msg.Message, "JOIN "))
			return
		}

		if strings.HasPrefix(msg.Message, "PART ") {
			b.Part(strings.TrimPrefix(msg.Message, "PART "))
			return
		}
	}

	b.listener.HandlePrivateMessage(msg)
}

func (b *Bot) joinBotConfigChannels() {
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
