package store

import (
	"encoding/json"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis() *Redis {
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

	return &Redis{
		Client: client,
	}
}
func (s *Redis) SubscribeSpeakerMessage() *redis.PubSub {
	return s.Client.Subscribe("SPEAKERMESSAGE")
}

type SpeakerMessage struct {
	Channel string
	Message string
}

func (s *Redis) PublishSpeakerMessage(channel, message string) {
	data, err := json.Marshal(&SpeakerMessage{channel, message})
	if err != nil {
		log.Error(err)
		return
	}

	s.Client.Publish("SPEAKERMESSAGE", data)
}

func (s *Redis) PublishPrivateMessage(message string) {
	s.Client.Publish("PRIVMSG", message)
}

func (s *Redis) SubscribePrivateMessages() *redis.PubSub {
	return s.Client.Subscribe("PRIVMSG")
}

const (
	IngesterMsgJoin = "join"
	IngesterMsgPart = "part"
)

type IngesterMessage struct {
	Type     string
	Argument string
}

func (s *Redis) PublishIngesterMessage(messageType string, argument string) {
	data, err := json.Marshal(&IngesterMessage{Type: messageType, Argument: argument})
	if err != nil {
		log.Error(err)
		return
	}

	s.Client.Publish("INGESTERMESSAGE", data)
}

func (s *Redis) SubscribeIngesterMessage() *redis.PubSub {
	return s.Client.Subscribe("PRIVMSG")
}
