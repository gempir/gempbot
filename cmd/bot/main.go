package main

import (
	"strings"
	"time"

	"github.com/gempir/gempbot/cmd/bot/commander"
	"github.com/gempir/gempbot/pkg/chat"
	"github.com/gempir/gempbot/pkg/config"
	"github.com/gempir/gempbot/pkg/helix"
	"github.com/gempir/gempbot/pkg/log"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/go-twitch-irc/v2"
)

// Bot basic logging bot
type Bot struct {
	startTime   time.Time
	cfg         *config.Config
	db          *store.Database
	helixClient *helix.Client
	listener    *commander.Listener
	Done        chan bool
	chatClient  *chat.ChatClient
}

func main() {
	cfg := config.FromEnv()

	db := store.NewDatabase(cfg)

	helixClient := helix.NewClient(cfg, db)
	go helixClient.StartRefreshTokenRoutine()

	bot := NewBot(cfg, db, helixClient)
	go bot.Connect()

	<-bot.Done
}

func NewBot(cfg *config.Config, db *store.Database, helixClient *helix.Client) *Bot {
	chatClient := chat.NewClient(cfg)

	handler := commander.NewHandler(cfg, helixClient, db, chatClient.Say)

	listener := commander.NewListener(db, handler, chatClient.Say)
	listener.RegisterDefaultCommands()

	return &Bot{
		Done:        make(chan bool),
		chatClient:  chatClient,
		cfg:         cfg,
		db:          db,
		listener:    listener,
		helixClient: helixClient,
	}
}

func (b *Bot) Connect() {
	b.startTime = time.Now()
	b.chatClient.SetOnPrivateMessage(b.handlePrivateMessage)
	go b.chatClient.Connect()

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("joining as anonymous")
	} else {
		log.Info("joining as user " + b.cfg.Username)
	}
	go b.chatClient.Join(b.cfg.Username)
	go b.joinBotConfigChannels()
}

func (b *Bot) handlePrivateMessage(msg twitch.PrivateMessage) {
	sysMessage := msg.Channel == b.cfg.Username && msg.User.Name == b.cfg.Username && strings.Contains(msg.Message, b.cfg.Environment)
	if sysMessage {
		log.Infof("sysMessage: %s", msg.Message)
		if strings.HasPrefix(msg.Message, "JOIN "+b.cfg.Environment+" ") {
			b.chatClient.Join(strings.TrimPrefix(msg.Message, "JOIN "+b.cfg.Environment+" "))
			return
		}

		if strings.HasPrefix(msg.Message, "PART "+b.cfg.Environment+" ") {
			b.chatClient.Part(strings.TrimPrefix(msg.Message, "PART "+b.cfg.Environment+" "))
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

	log.Infof("joining %d channels", len(users))
	for _, user := range users {
		time.Sleep(time.Millisecond * 1000)
		b.chatClient.Join(user.Login)
	}
}
