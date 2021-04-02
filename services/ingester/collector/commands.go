package collector

import (
	"fmt"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/pkg/humanize"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) handlePrivateMessage(message twitch.PrivateMessage) {
	b.active.mutex.Lock()
	b.active.m[strings.ToLower(message.Channel)] = true
	b.active.mutex.Unlock()
	b.store.PublishPrivateMessage(message.Raw)
	if message.User.Name == b.cfg.Admin {
		if strings.HasPrefix(message.Message, "!spamchamp status") {
			uptime := humanize.TimeSince(b.startTime)
			b.scaler.Say(message.Channel, message.User.DisplayName+", uptime: "+uptime)
		}
		if strings.HasPrefix(message.Message, "!spamchamp join ") {
			b.handleJoin(message)
		}
		if strings.HasPrefix(message.Message, "!spamchamp top") {
			b.LoadTopChannelsAndJoin()
		}
	}

	if val, ok := message.User.Badges["partner"]; ok && val == 1 {
		if _, ok := b.joined.m[message.User.Name]; !ok {
			log.Infof("Found partner, joining channel: %s", message.User.Name)
			b.scaler.Join(message.User.Name)
			b.store.AddChannels(message.User.ID)
			b.joined.mutex.Lock()
			b.joined.m[message.User.Name] = true
			b.joined.mutex.Unlock()
		}
	}
}

func (b *Bot) handleJoin(message twitch.PrivateMessage) {
	input := strings.TrimPrefix(message.Message, "!spamchamp join ")

	users, err := b.helixClient.GetUsersByUsernames(strings.Split(input, ","))
	if err != nil {
		log.Error(err)
		b.scaler.Say(message.Channel, message.User.DisplayName+", something went wrong requesting the userids")
	}

	ids := []string{}
	for _, user := range users {
		ids = append(ids, user.ID)
		log.Infof("[collector] joining %s", user.Login)
		b.scaler.Join(user.Login)
	}
	b.store.AddChannels(ids...)
	b.joinStoreChannels()
	b.scaler.Say(message.Channel, fmt.Sprintf("%s, added channels: %v", message.User.DisplayName, ids))
}
