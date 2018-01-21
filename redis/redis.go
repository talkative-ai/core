package redis

import (
	"os"

	"github.com/go-redis/redis"
)

var Instance *redis.Client

// ConnectRedis dials and returns a Redis client to query against, or abort with error
func ConnectRedis() (*redis.Client, error) {

	if os.Getenv("REDIS_ADDR") == "" {
		os.Setenv("REDIS_ADDR", "127.0.0.1:6379")
		os.Setenv("REDIS_PASSWORD", "")
	}

	Instance = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := Instance.Ping().Result()
	if err != nil {
		return nil, err
	}

	return Instance, nil
}
