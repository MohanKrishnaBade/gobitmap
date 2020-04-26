package routes

import (
	"github.com/go-redis/redis/v7"
	"github.com/gobitmap/redisStorage"
)

var redisClient *redis.Client

func init() {
	redisClient = redisStorage.RedisConnect()
}

func Get(key string) (string, error) {
	return redisClient.Get(key).Result()
}

func Set(key, value string) error {

	return redisClient.Set(key, value, 0).Err()
}

func CreateUnless(key, value string) error {
	_, err := Get(key)
	if err == redis.Nil || err != nil {
		return Set(key, value)
	}
	return nil
}
