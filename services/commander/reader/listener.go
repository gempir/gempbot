package reader

import (
	"strings"

	"github.com/gempir/bitraft/pkg/store"
	"github.com/gempir/bitraft/pkg/tmi"
	"github.com/gempir/bitraft/services/commander/predictions"
	"github.com/gempir/go-twitch-irc/v2"
)

type Listener struct {
	db                 *store.Database
	redis              *store.Redis
	predictionsHandler *predictions.Handler
	commands           map[string]func(twitch.PrivateMessage)
}

func NewListener(db *store.Database, redis *store.Redis, predictionsHandler *predictions.Handler) *Listener {

	return &Listener{
		db:                 db,
		redis:              redis,
		predictionsHandler: predictionsHandler,
		commands:           map[string]func(twitch.PrivateMessage){},
	}
}

func (l *Listener) RegisterDefaultCommands() {
	l.commands["!prediction"] = l.handlePrediction
}

func (l *Listener) StartListener() {
	topic := l.redis.SubscribePrivateMessages()
	channel := topic.Channel()
	for msg := range channel {
		parsedMsg := twitch.ParseMessage(msg.Payload).(*twitch.PrivateMessage)

		l.handleMessage(*parsedMsg)
	}
}

func (l *Listener) handleMessage(msg twitch.PrivateMessage) {
	if !strings.HasPrefix(msg.Message, "!") {
		return
	}

	groups := strings.Split(msg.Message, " ")

	if cmd, ok := l.commands[groups[0]]; ok {
		cmd(msg)
	}
}

func (l *Listener) handlePrediction(msg twitch.PrivateMessage) {
	perm := l.db.GetPermission(msg.User.ID, msg.RoomID)
	if !perm.Prediction && !tmi.IsModerator(msg.User) && !tmi.IsBroadcaster(msg.User) {
		return
	}

	l.predictionsHandler.HandleMessage(msg)
}
