package stats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/gempir/spamchamp/bot/api"
	"github.com/gempir/spamchamp/bot/store"
	"github.com/paulbellamy/ratecounter"
	log "github.com/sirupsen/logrus"
)

var (
	stats          = map[string]stat{}
	joinedChannels = 0
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
		if message.ID == "28b511cc-43b3-44b7-a605-230aadbb2f9b" {
			log.Info(message.Message)
			var err error
			joinedChannels, err = strconv.Atoi(message.Message)
			if err != nil {
				log.Errorf("Failed to parse relaybroker message: %s", err.Error())
			}
			continue
		}
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
			Records: []api.Record{},
		}

		msgps := api.Record{
			Title:  "Current messages/s",
			Scores: []api.Score{},
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

			msgps.Scores = append(msgps.Scores, api.Score{ID: channelID, Score: float64(rate)})
		}

		scores := []api.Score{}
		for _, z := range b.store.GetMsgpsScores() {
			scores = append(scores, api.Score{ID: fmt.Sprintf("%v", z.Member), Score: z.Score})
		}

		maxLen := 10
		if len(msgps.Scores) < 10 {
			maxLen = len(msgps.Scores)
		}
		msgps.Scores = msgps.GetScoresSorted()[0:maxLen]
		message.Records = append(message.Records, msgps)

		message.Records = append(message.Records, api.Record{
			Title:  "Record messages/s",
			Scores: scores,
		})

		message.JoinedChannels = joinedChannels

		b.broadcastQueue <- message
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
