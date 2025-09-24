package cache

import (
	"context"
	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/wb-go/wbf/redis"
	"os"
	"time"
)

type Repository interface {
	GetAllNotifications() ([]models.Notification, error)
}

type Redis struct {
	client *redis.Client
}

func New(host, port string) *Redis {
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.New(
		host+port,
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

func (r *Redis) Set(key string, val interface{}) error {
	return r.client.SetEX(context.Background(), key, val, 24*time.Hour).Err()
}

//func (r *Redis) LoadAllNotifications(r *repository.Repository) error {
//	// TODO: дописать реализацию
//	//notifications, err := r.GetAllNotifications()
//	//if err != nil {
//	//	zlog.Logger.Fatal(err).Msg("error get all notifications")
//	//}
//	return nil
//}
