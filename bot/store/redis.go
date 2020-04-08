package store

import (
	"fmt"

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

func (s *Store) UpdateMsgps(channelID string, msgps int64) {
	s.redis.ZAdd("channel:msgps", &redis.Z{Score: float64(msgps), Member: channelID})
}

func (s *Store) GetMsgps(channelID string) float64 {
	score, err := s.redis.ZScore("channel:msgps", channelID).Result()
	if err != nil {
		return 0
	}

	return score
}

func (s *Store) GetMsgpsScores() []redis.Z {
	scores, err := s.redis.ZRevRangeWithScores("channel:msgps", 0, 9).Result()
	if err != nil {
		return []redis.Z{}
	}

	return scores
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

type cloudWord struct {
	Word string
	Add  int64
}

func (s *Store) AddToWordcloud(words ...string) {
	for _, word := range words {
		_, err := s.redis.ZIncrBy("wordcloud", 1, word).Result()
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

func (s *Store) GetEntireWordcloud() map[string]float64 {
	words, err := s.redis.ZRangeWithScores("wordcloud", 0, -1).Result()
	if err != nil {
		log.Error(err)
		return map[string]float64{}
	}

	wordcloudWords := map[string]float64{}
	for _, value := range words {
		word := fmt.Sprintf("%v", value.Member)

		wordcloudWords[word] = value.Score
		if value.Score == 0 {
			s.redis.ZRem("wordcloud")
		}
	}

	return wordcloudWords
}

func (s *Store) GetTopWords() map[string]float64 {
	words, err := s.redis.ZRevRangeWithScores("wordcloud", 0, 49).Result()
	if err != nil {
		log.Error(err)
		return map[string]float64{}
	}

	wordcloudWords := map[string]float64{}
	for _, value := range words {
		word := fmt.Sprintf("%v", value.Member)
		wordcloudWords[word] = value.Score
	}

	return wordcloudWords
}

func (s *Store) TickDownAll() {
	_, err := s.redis.Eval(`
	local zsetMembers = redis.call('zrange', KEYS[1], '0', '-1') 
	for k,member in pairs(zsetMembers) do 
		redis.call('zincrby', KEYS[1], -1, member) 
	end`,
		[]string{"wordcloud"}).Result()
	if err != nil && err.Error() != "redis: nil" {
		log.Error(err.Error())
	}
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
