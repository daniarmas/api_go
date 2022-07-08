package rdb

import (
	"fmt"

	"github.com/daniarmas/api_go/config"
	"github.com/go-redis/redis/v9"
)

func New(cfg *config.Config) *redis.Client {
	rdbAddres := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr: rdbAddres,
		// Password: config.RedisPassword, // no password set
		// DB:       config.RedisDb,       // use default DB
	})
	return rdb
}
