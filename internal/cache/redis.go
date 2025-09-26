package cache

import (
	"context"
	"github.com/K1la/delayed-notifier/internal/config"
	"os"
	"time"

	"github.com/wb-go/wbf/redis"
)

type Redis struct {
	client *redis.Client
}

func New(cfg config.RedisCfg) *Redis {
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.New(
		cfg.Host+":"+cfg.Port,
		password,
		0,
	)

	return &Redis{
		client: client,
	}
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.Background(), key)
}

func (r *Redis) Set(key string, val any) error {
	return r.client.SetEX(context.Background(), key, val, 24*time.Hour).Err()
}
