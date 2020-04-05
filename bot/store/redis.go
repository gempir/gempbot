package store

import (
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	redis *redis.Client
}

func NewStore() Store {
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

	return Store{
		redis: client,
	}
}
