package redis

import (
	"context"
	"file-server/setting"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	DB  *redis.Client
	Ctx = context.Background()
)

func Init(cfg *setting.RedisConfig) {
	DB = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Password: "",
		DB:       cfg.Db,
		PoolSize: cfg.PoolSize,
	})
}
