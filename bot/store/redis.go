package store

import (
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	redis *redis.Client
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
		redis: client,
	}
}

func (s *Store) AddChannels(channelIDs ...string) {
	for _, id := range channelIDs {
		_, err := s.redis.HSet("channels", id, "1").Result()
		if err != nil {
			log.Error(err)
			continue
		}
	}
	log.Infof("[store] added %v", channelIDs)
}

func (s *Store) RemoveChannel(channelID string) {
	_, err := s.redis.HDel("channels", channelID).Result()
	if err != nil {
		log.Error(err)
	}
	log.Infof("[store] removed %s", channelID)
}

func (s *Store) GetAllChannels() []string {
	channels, err := s.redis.HKeys("channels").Result()
	if err != nil {
		log.Error(err)
		return []string{}
	}

	return channels
}
