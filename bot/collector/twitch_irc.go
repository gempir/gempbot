package collector

import (
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/config"
	"github.com/gempir/spamchamp/bot/helix"
	log "github.com/sirupsen/logrus"
)

// Bot basic logging bot
type Bot struct {
	messageQueue chan twitch.PrivateMessage
	startTime    time.Time
	cfg          *config.Config
	helixClient  *helix.Client
	twitchClient *twitch.Client
	channels     map[string]helix.UserData
}

// NewBot create new bot instance
func NewBot(cfg *config.Config, helixClient *helix.Client, messageQueue chan twitch.PrivateMessage) *Bot {
	channels, err := helixClient.GetUsersByUserIds(cfg.Channels)
	if err != nil {
		log.Fatalf("[collector] failed to load configured channels %s", err.Error())
	}

	return &Bot{
		messageQueue: messageQueue,
		cfg:          cfg,
		helixClient:  helixClient,
		channels:     channels,
	}
}

// Connect startup the logger and bot
func (b *Bot) Connect() {
	b.startTime = time.Now()
	b.twitchClient = twitch.NewClient(b.cfg.Username, "oauth:"+b.cfg.OAuth)

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

	log.Fatal(b.twitchClient.Connect())
}

func (b *Bot) initialJoins() {
	for _, channel := range b.channels {
		log.Info("[collector] joining " + channel.Login)
		b.twitchClient.Join(channel.Login)
	}
}
