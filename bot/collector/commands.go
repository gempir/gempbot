package collector

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/humanize"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handlePrivateMessage(message twitch.PrivateMessage) {
	if message.User.Name == b.cfg.Admin {
		if strings.HasPrefix(message.Message, "!spamchamp status") {
			uptime := humanize.TimeSince(b.startTime)
			b.twitchClient.Say(message.Channel, message.User.DisplayName+", uptime: "+uptime)
		}
		if strings.HasPrefix(message.Message, "!spamchamp join ") {
			b.handleJoin(message)
		}
	}
}

func (b *Bot) handleJoin(message twitch.PrivateMessage) {
	input := strings.TrimPrefix(message.Message, "!spamchamp join ")

	users, err := b.helixClient.GetUsersByUsernames(strings.Split(input, ","))
	if err != nil {
		log.Error(err)
		b.twitchClient.Say(message.Channel, message.User.DisplayName+", something went wrong requesting the userids")
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
		log.Infof("[bot] joining %s", user.Login)
		b.twitchClient.Join(user.Login)
	}
	b.cfg.AddChannels(ids...)
	b.twitchClient.Say(message.Channel, fmt.Sprintf("%s, added channels: %v", message.User.DisplayName, ids))
}
