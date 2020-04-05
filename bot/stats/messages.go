package stats

import (
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/api"
	"github.com/paulbellamy/ratecounter"
	log "github.com/sirupsen/logrus"
)

var (
	stats = map[string]stat{}
)

type Broadcaster struct {
	messageQueue   chan twitch.PrivateMessage
	broadcastQueue chan api.BroadcastMessage
}

func NewBroadcaster(messageQueue chan twitch.PrivateMessage, broadcastQueue chan api.BroadcastMessage) Broadcaster {
	return Broadcaster{
		messageQueue:   messageQueue,
		broadcastQueue: broadcastQueue,
	}
}

func (b *Broadcaster) Start() {
	log.Info("[collector] starting stats collector")

	go b.startTicker()

	for message := range b.messageQueue {
		if _, ok := stats[message.RoomID]; !ok {
			stats[message.RoomID] = newStat(message.Channel)
		}

		stats[message.RoomID].messages.Incr(1)
	}
}

func (b *Broadcaster) startTicker() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		message := api.BroadcastMessage{
			ChannelStats: []api.ChannelStat{},
		}

		for channelID, stat := range stats {
			rate := stat.messages.Rate()
			if rate == 0 {
				continue
			}

			message.ChannelStats = append(message.ChannelStats, api.ChannelStat{
				ID:    channelID,
				Msgps: rate,
			})
		}

		if len(stats) > 0 {
			b.broadcastQueue <- message
		}
	}
}

type stat struct {
	channelName  string
	messages     *ratecounter.RateCounter
	messageCount int
}

func newStat(channelName string) stat {
	return stat{
		channelName: channelName,
		messages:    ratecounter.NewRateCounter(1 * time.Second),
	}
}
