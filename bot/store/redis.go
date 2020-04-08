package store

import (
	"strconv"

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
		_, err := s.redis.HIncrBy("wordcloud", word, 1).Result()
		if err != nil {
			log.Error(err)
			continue
		}
	}
}

func (s *Store) GetEntireWordcloud() map[string]int {
	words, err := s.redis.HGetAll("wordcloud").Result()
	if err != nil {
		log.Error(err)
		return map[string]int{}
	}

	wordcloudWords := map[string]int{}
	for word, value := range words {
		res, err := strconv.Atoi(value)
		if err != nil {
			log.Error(err)
			continue
		}

		wordcloudWords[word] = res
		if res == 0 {
			s.redis.HDel("wordcloud", word)
		}
	}

	return wordcloudWords
}

func (s *Store) TickDownWord(word string) {
	_, err := s.redis.HIncrBy("wordcloud", word, -1).Result()
	if err != nil {
		log.Error(err)
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
