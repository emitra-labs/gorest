package store

import (
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisLock *redislock.Client

func Open(redisURL string) error {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}

	RedisClient = redis.NewClient(opts)
	RedisLock = redislock.New(RedisClient)

	return nil
}

func Close() error {
	return RedisClient.Close()
}
