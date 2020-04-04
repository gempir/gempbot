package stats

import (
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/api"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	stats = map[string]stat{}
)


type Broadcaster struct {
	messageQueue chan twitch.PrivateMessage
	broadcastQueue chan api.BroadcastMessage
}

func NewBroadcaster(messageQueue chan twitch.PrivateMessage, broadcastQueue chan api.BroadcastMessage) Broadcaster {
	return Broadcaster{
		messageQueue: messageQueue,
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

		stats[message.RoomID].messages[0]++
	}
}

func (b *Broadcaster) startTicker() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		message := api.BroadcastMessage{
			ChannelStats: []api.ChannelStat{},
		}

		for channelID, stat := range stats {
			ratebucket := stat.messages
			total := 0
			for i := len(ratebucket) - 1; i >= 1; i-- {
				total += ratebucket[i]
				ratebucket[i] = ratebucket[i-1]
			}
			total += ratebucket[0]
			ratebucket[0] = 0

			message.ChannelStats = append(message.ChannelStats, api.ChannelStat{
				ID:    channelID,
				Msgps: total / 10,
				Msgpm: total * 6,
			})
		}

		if len(stats) > 0 {
			b.broadcastQueue <- message
		}
	}
}

type stat struct {
	channelName string
	messages    []int
}

func newStat(channelName string) stat {
	return stat{
		channelName: channelName,
		messages:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}