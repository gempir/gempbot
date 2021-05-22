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
	store  *store.Store
	client *twitch.Client
}

func NewBot(cfg *config.Config, store *store.Store) *Bot {
	return &Bot{
		cfg:    cfg,
		store:  store,
		client: nil,
	}
}

func (b *Bot) Connect() {
	b.client = twitch.NewClient(b.cfg.Username, b.cfg.OAuth)

	err := b.client.Connect()
	if err != nil {
		panic(err)
	}
}

func (b *Bot) Say(channel, message string) {
	log.Infof("[%s]: %s", channel, message)
	b.client.Say(channel, humanize.CharLimiter(message, 500))
}
