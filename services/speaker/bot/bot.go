package bot

import (
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/pkg/config"
	"github.com/gempir/spamchamp/pkg/store"
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
	b.client.Say(channel, message)
}
