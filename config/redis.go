package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		panic("redis connection failed:" + err.Error())
	}
}
