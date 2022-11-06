package commander

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gempir/gempbot/internal/chat/tmi"
	"github.com/gempir/gempbot/internal/dto"
	"github.com/gempir/gempbot/internal/humanize"
	"github.com/gempir/gempbot/internal/store"
	"github.com/gempir/go-twitch-irc/v3"
)

type Listener struct {
	startTime          time.Time
	db                 *store.Database
	predictionsHandler *Handler
	commands           map[string]func(dto.CommandPayload)
	chatSay            func(channel, message string)
}

var (
	commandRegex = regexp.MustCompile(`^\!(\w+)\ ?`)
)

func NewListener(db *store.Database, predictionsHandler *Handler, chatSay func(channel, message string)) *Listener {
	return &Listener{
		startTime:          time.Now(),
		db:                 db,
		predictionsHandler: predictionsHandler,
		commands:           map[string]func(dto.CommandPayload){},
		chatSay:            chatSay,
	}
}

func (l *Listener) RegisterCommand(command string, handler func(dto.CommandPayload)) {
	l.commands[command] = handler
}

func (l *Listener) RegisterDefaultCommands() {
	l.commands[dto.CmdNameStatus] = l.handleStatus
	l.commands[dto.CmdNamePrediction] = l.handlePrediction
	l.commands[dto.CmdNameOutcome] = l.handlePrediction
}

func (l *Listener) HandlePrivateMessage(msg twitch.PrivateMessage) {
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

func (l *Listener) hasElevatedPermissions(payload dto.CommandPayload) bool {
	return tmi.IsModerator(payload.Msg.User) || tmi.IsBroadcaster(payload.Msg.User) || l.db.GetChannelUserPermissions(payload.Msg.User.ID, payload.Msg.RoomID).Prediction
}

func (l *Listener) handlePrediction(payload dto.CommandPayload) {
	if !l.hasElevatedPermissions(payload) {
		return
	}

	l.predictionsHandler.HandleCommand(payload)
}

func (l *Listener) handleStatus(payload dto.CommandPayload) {
	if !l.hasElevatedPermissions(payload) {
		return
	}

	uptime := humanize.TimeSince(l.startTime)
	l.chatSay(payload.Msg.Channel, fmt.Sprintf("@%s, uptime: %s", payload.Msg.User.DisplayName, uptime))
}
