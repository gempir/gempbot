package stats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/pkg/store"
	"github.com/gempir/spamchamp/services/api/server"
	"github.com/paulbellamy/ratecounter"
	log "github.com/sirupsen/logrus"
)

var (
	stats          = map[string]stat{}
	joinedChannels = 0
	activeChannels = 0
)

type Broadcaster struct {
	broadcastQueue chan server.BroadcastMessage
	store          *store.Store
}

func NewBroadcaster(broadcastQueue chan server.BroadcastMessage, store *store.Store) Broadcaster {
	return Broadcaster{
		broadcastQueue: broadcastQueue,
		store:          store,
	}
}

func (b *Broadcaster) Start() {
	log.Info("[stats] starting stats broadcaster")

	go b.startTicker()

	go b.monitorJoinedChannels()
	go b.monitorActiveChannels()

	topic := b.store.SubscribePrivateMessages()
	channel := topic.Channel()
	for msg := range channel {
		message := twitch.ParseMessage(msg.Payload).(*twitch.PrivateMessage)

		if _, ok := stats[message.RoomID]; !ok {
			stats[message.RoomID] = newStat(message.Channel)
		}

		stats[message.RoomID].messages.Incr(1)
	}
}

func (b *Broadcaster) monitorJoinedChannels() {
	topic := b.store.SubscribeJoinedChannels()
	channel := topic.Channel()
	for msg := range channel {

		var err error
		joinedChannels, err = strconv.Atoi(msg.Payload)
		if err != nil {
			log.Errorf("Failed to parse joined channels message: %s", err.Error())
		}
		continue
	}
}

func (b *Broadcaster) monitorActiveChannels() {
	topic := b.store.SubscribeActiveChannels()
	channel := topic.Channel()
	for msg := range channel {

		var err error
		activeChannels, err = strconv.Atoi(msg.Payload)
		if err != nil {
			log.Errorf("Failed to parse active channels message: %s", err.Error())
		}
		continue
	}
}

func (b *Broadcaster) startTicker() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		message := server.BroadcastMessage{
			Records: []server.Record{},
		}

		msgps := server.Record{
			Title:  "Current messages/s",
			Scores: []server.Score{},
		}

		for channelID, stat := range stats {
			rate := stat.messages.Rate() / 3
			if rate == 0 {
				continue
			}

			score := b.store.GetMsgps(channelID)
			if float64(rate) > score {
				b.store.UpdateMsgps(channelID, rate)
			}

			msgps.Scores = append(msgps.Scores, server.Score{ID: channelID, Score: float64(rate)})
		}

		scores := []server.Score{}
		for _, z := range b.store.GetMsgpsScores() {
			scores = append(scores, server.Score{ID: fmt.Sprintf("%v", z.Member), Score: z.Score})
		}

		maxLen := 10
		if len(msgps.Scores) < 10 {
			maxLen = len(msgps.Scores)
		}
		msgps.Scores = msgps.GetScoresSorted()[0:maxLen]
		message.Records = append(message.Records, msgps)

		message.Records = append(message.Records, server.Record{
			Title:  "Record messages/s",
			Scores: scores,
		})

		message.JoinedChannels = joinedChannels
		message.ActiveChannels = activeChannels

		b.broadcastQueue <- message
	}
}

type stat struct {
	channelName string
	messages    *ratecounter.RateCounter
}

func newStat(channelName string) stat {
	return stat{
		channelName: channelName,
		messages:    ratecounter.NewRateCounter(time.Second * 3),
	}
}
