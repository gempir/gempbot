package bot

import (
	"github.com/gempir/bitraft/pkg/config"
	"github.com/gempir/bitraft/pkg/humanize"
	"github.com/gempir/bitraft/pkg/store"
	"github.com/gempir/go-twitch-irc/v2"
	log "github.com/sirupsen/logrus"
)

type Bot struct {
	cfg    *config.Config
	store  *store.Redis
	client *twitch.Client
	db     *store.Database
}

func NewBot(cfg *config.Config, store *store.Redis, db *store.Database) *Bot {
	return &Bot{
		cfg:    cfg,
		store:  store,
		client: nil,
		db:     db,
	}
}

func (b *Bot) Connect() {
	b.client = twitch.NewClient(b.cfg.Username, b.cfg.OAuth)

	err := b.client.Connect()
	if err != nil {
		panic(err)
	}
}

func (b *Bot) Say(userid, channel, message string) {
	cfg, err := b.db.GetBotConfig(userid)
	if err != nil || !cfg.JoinBot {
		log.Warnf("[%s]: %s - no permission to Say %s", channel, message, err)
		return
	}

	log.Infof("[%s]: %s", channel, message)
	b.client.Say(channel, humanize.CharLimiter(message, 500))
}
