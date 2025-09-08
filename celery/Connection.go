package celery

import (
	redispool "go-restapi-gin/redispool"

	"github.com/gocelery/gocelery"
)

func Connect() (*gocelery.CeleryClient, error) {
	redis, _ := redispool.Connect()
	cli, _ := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(&redis),
		&gocelery.RedisCeleryBackend{Pool: &redis},
		5, // number of workers
	)

	return cli, nil
}
