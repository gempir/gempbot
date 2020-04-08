package stats

import (
	"fmt"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/api"
	"github.com/gempir/spamchamp/bot/store"
	"github.com/paulbellamy/ratecounter"
	log "github.com/sirupsen/logrus"
)

var (
	stats          = map[string]stat{}
	activeChannels = map[string]bool{}
)

type Broadcaster struct {
	messageQueue   chan twitch.PrivateMessage
	broadcastQueue chan api.BroadcastMessage
	store          *store.Store
}

func NewBroadcaster(messageQueue chan twitch.PrivateMessage, broadcastQueue chan api.BroadcastMessage, store *store.Store) Broadcaster {
	return Broadcaster{
		messageQueue:   messageQueue,
		broadcastQueue: broadcastQueue,
		store:          store,
	}
}

func (b *Broadcaster) Start() {
	log.Info("[stats] starting stats collector")

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
			Records:      []api.Record{},
		}

		for channelID, stat := range stats {
			rate := stat.messages.Rate() / 3
			if rate == 0 {
				continue
			}
			activeChannels[channelID] = true

			score := b.store.GetMsgps(channelID)
			if float64(rate) > score {
				b.store.UpdateMsgps(channelID, rate)
			}

			message.ChannelStats = append(message.ChannelStats, api.ChannelStat{
				ID:    channelID,
				Msgps: rate,
			})
		}

		scores := []api.Score{}
		for _, z := range b.store.GetMsgpsScores() {
			scores = append(scores, api.Score{ID: fmt.Sprintf("%v", z.Member), Score: z.Score})
		}

		message.Records = append(message.Records, api.Record{
			Title:  "Records: messages/s",
			Scores: scores,
		})

		message.ActiveChannels = len(activeChannels)

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
		messages:    ratecounter.NewRateCounter(time.Second * 3),
	}
}
