package store

import (
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
	log.Info("[store] connection to redis established")

	return &Store{
		Client: client,
	}
}

func (s *Store) UpdateMsgps(channelID string, msgps int64) {
	s.Client.ZAdd("channel:msgps", &redis.Z{Score: float64(msgps), Member: channelID})
}

func (s *Store) GetMsgps(channelID string) float64 {
	score, err := s.Client.ZScore("channel:msgps", channelID).Result()
	if err != nil {
		return 0
	}

	return score
}

func (s *Store) GetMsgpsScores() []redis.Z {
	scores, err := s.Client.ZRevRangeWithScores("channel:msgps", 0, 9).Result()
	if err != nil {
		return []redis.Z{}
	}

	return scores
}

func (s *Store) AddChannels(channelIDs ...string) {
	for _, id := range channelIDs {
		_, err := s.Client.HSet("channels", id, "1").Result()
		if err != nil {
			log.Error(err)
			continue
		}
	}
	log.Infof("[store] added %v", channelIDs)
}

func (s *Store) PublishPrivateMessage(message string) {
	s.Client.Publish("PRIVMSG", message)
}

func (s *Store) SubscribePrivateMessages() *redis.PubSub {
	return s.Client.Subscribe("PRIVMSG")
}

func (s *Store) PublishJoinedChannels(count int) {
	s.Client.Publish("JOINEDCHANNELS", count)
}

func (s *Store) SubscribeJoinedChannels() *redis.PubSub {
	return s.Client.Subscribe("JOINEDCHANNELS")
}

func (s *Store) PublishActiveChannels(count int) {
	s.Client.Publish("ACTIVECHANNELS", count)
}

func (s *Store) SubscribeActiveChannels() *redis.PubSub {
	return s.Client.Subscribe("ACTIVECHANNELS")
}

func (s *Store) RemoveChannel(channelID string) {
	_, err := s.Client.HDel("channels", channelID).Result()
	if err != nil {
		log.Error(err)
	}
	log.Infof("[store] removed %s", channelID)
}

func (s *Store) GetAllChannels() []string {
	channels, err := s.Client.HKeys("channels").Result()
	if err != nil {
		log.Error(err)
		return []string{}
	}

	return channels
}
