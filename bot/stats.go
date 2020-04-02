package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	stats = map[string]stat{}
)

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

func startStatsCollector() {
	log.Info("[collector] starting stats collector")

	go startTicker()

	for message := range messageQueue {
		if _, ok := stats[message.RoomID]; !ok {
			stats[message.RoomID] = newStat(message.Channel)
		}

		stats[message.RoomID].messages[0]++
	}
}

func startTicker() {
	ticker := time.NewTicker(1 * time.Second)

	for range ticker.C {
		message := socketMessage{
			Channels: make(map[string]frontendStats),
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

			message.Channels[channelID] = frontendStats{
				ChannelName:       stat.channelName,
				MessagesPerSecond: total / 10,
			}
		}

		broadcast <- message
	}
}
