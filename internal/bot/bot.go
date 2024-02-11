package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/gempir/gempbot/internal/chat"
	"github.com/gempir/gempbot/internal/config"
	"github.com/gempir/gempbot/internal/helixclient"
	"github.com/gempir/gempbot/internal/log"
	"github.com/gempir/gempbot/internal/store"
	"github.com/nicklaw5/helix/v2"
)

const gempbotUserID = "99659894"

// Bot basic logging bot
type Bot struct {
	startTime   time.Time
	cfg         *config.Config
	db          *store.Database
	helixClient helixclient.Client
	Done        chan bool
	ChatClient  *chat.ChatClient
}

func NewBot(cfg *config.Config, db *store.Database, helixClient helixclient.Client) *Bot {
	chatClient := chat.NewClient(cfg, helixClient)

	return &Bot{
		Done:        make(chan bool),
		ChatClient:  chatClient,
		cfg:         cfg,
		db:          db,
		helixClient: helixClient,
	}
}

func (b *Bot) Send(channelID string, message string) {
	fmt.Println("sending message", channelID, message, gempbotUserID)
	resp, err := b.helixClient.SendChatMessage(&helix.SendChatMessageParams{BroadcasterID: channelID, Message: message, SenderID: gempbotUserID})
	if err != nil {
		log.Error("Failure sending message", err, resp)
	}
	fmt.Println(resp)
}

func (b *Bot) Join(channel string) {
	go b.ChatClient.Join(channel)
}

func (b *Bot) Part(channel string) {
	go b.ChatClient.Part(channel)
}

func (b *Bot) Connect() {
	b.startTime = time.Now()
	go b.ChatClient.Connect(b.joinBotConfigChannels)

	if strings.HasPrefix(b.cfg.Username, "justinfan") {
		log.Info("joining as anonymous")
	} else {
		log.Info("joining as user " + b.cfg.Username)
	}
	go b.ChatClient.Join(b.cfg.Username)
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

	b.ChatClient.WaitForConnect()
	log.Infof("joining %d channels", len(users))
	for _, user := range users {
		b.ChatClient.Join(user.Login)
	}
}
