package reader

import (
	"regexp"
	"strings"

	"github.com/gempir/gempbot/cmd/bot/commander/predictions"
	"github.com/gempir/gempbot/pkg/dto"
	"github.com/gempir/gempbot/pkg/store"
	"github.com/gempir/gempbot/pkg/tmi"
	"github.com/gempir/go-twitch-irc/v2"
)

type Listener struct {
	db                 *store.Database
	redis              *store.Redis
	predictionsHandler *predictions.Handler
	commands           map[string]func(dto.CommandPayload)
}

var (
	commandRegex = regexp.MustCompile(`^\!(\w+)\ ?`)
)

func NewListener(db *store.Database, redis *store.Redis, predictionsHandler *predictions.Handler) *Listener {
	return &Listener{
		db:                 db,
		redis:              redis,
		predictionsHandler: predictionsHandler,
		commands:           map[string]func(dto.CommandPayload){},
	}
}

func (l *Listener) RegisterDefaultCommands() {
	l.commands[dto.CmdNamePrediction] = l.handlePrediction
	l.commands[dto.CmdNameOutcome] = l.handlePrediction
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

	match := commandRegex.FindStringSubmatch(msg.Message)
	if len(match) < 2 {
		return
	}

	if cmd, ok := l.commands[match[1]]; ok {
		cmd(dto.CommandPayload{Msg: msg, Name: match[1], Query: strings.TrimSpace(strings.TrimPrefix(msg.Message, "!"+match[1]))})
	}
}

func (l *Listener) handlePrediction(payload dto.CommandPayload) {
	perm := l.db.GetChannelUserPermissions(payload.Msg.User.ID, payload.Msg.RoomID)
	if !perm.Prediction && !tmi.IsModerator(payload.Msg.User) && !tmi.IsBroadcaster(payload.Msg.User) {
		return
	}

	l.predictionsHandler.HandleCommand(payload)
}
