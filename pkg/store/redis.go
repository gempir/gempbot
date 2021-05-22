package store

import (
	"encoding/json"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	Client *redis.Client
}

func NewStore() *Store {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("connection to redis established")

	return &Store{
		Client: client,
	}
}
func (s *Store) SubscribeSpeakerMessage() *redis.PubSub {
	return s.Client.Subscribe("SPEAKERMESSAGE")
}

type SpeakerMessage struct {
	Channel string
	Message string
}

func (s *Store) PublishSpeakerMessage(channel, message string) {
	data, err := json.Marshal(&SpeakerMessage{channel, message})
	if err != nil {
		log.Error(err)
		return
	}

	s.Client.Publish("SPEAKERMESSAGE", data)
}
