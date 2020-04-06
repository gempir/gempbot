package collector

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/config"
	"github.com/gempir/spamchamp/bot/helix"
	"github.com/gempir/spamchamp/bot/store"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	messageQueue chan twitch.PrivateMessage
	startTime    time.Time
	cfg          *config.Config
	helixClient  *helix.Client
	twitchClient *twitch.Client
	store        *store.Store
	channels     map[string]helix.UserData
}

// NewBot create new bot instance
func NewBot(cfg *config.Config, helixClient *helix.Client, store *store.Store, messageQueue chan twitch.PrivateMessage) *Bot {
	channels, err := helixClient.GetUsersByUserIds(cfg.Channels)
	if err != nil {
		log.Fatalf("[collector] failed to load configured channels %s", err.Error())
	}

	return &Bot{
		messageQueue: messageQueue,
		cfg:          cfg,
		helixClient:  helixClient,
		store:        store,
		channels:     channels,
	}
}

// Connect startup the logger and bot
func (b *Bot) Connect() {
	b.startTime = time.Now()
	b.twitchClient = twitch.NewClient(b.cfg.Username, "oauth:"+b.cfg.OAuth)
	b.twitchClient.IrcAddress = "127.0.0.1:3333"
	b.twitchClient.TLS = false
	// b.twitchClient.SetIRCToken("LOGIN spamchampbot")

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("[collector] joining as anonymous")
	} else {
		log.Info("[collector] joining as user " + b.cfg.Username)
	}
	b.initialJoins()

	b.twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		b.messageQueue <- message
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

func (b *Bot) slowlyJoinStoreChannels() {
	go func() {
		channels, err := b.helixClient.GetUsersByUserIds(b.store.GetAllChannels())
		if err != nil {
			log.Error(err)
			return
		}

		for _, userData := range channels {
			b.twitchClient.Join(userData.Login)
			log.Debugf("[collector] slowly joined %s", userData.DisplayName)
			// time.Sleep(time.Second)
		}
	}()
}

func (b *Bot) LoadTopChannelsAndJoin() {
	log.Info("[collector] fetching top channels and joining")
	b.store.AddChannels(b.helixClient.GetTopChannels()...)

	b.slowlyJoinStoreChannels()
}
