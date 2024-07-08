package redisconn

import (
	"os"

	redis "github.com/redis/go-redis/v9"
)

func Connect() (*redis.Client, error) {

	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,
	}), nil
}
