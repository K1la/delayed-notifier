package service

import (
	"github.com/K1la/delayed-notifier/internal/cache"
	"github.com/K1la/delayed-notifier/internal/repository"
)

type NotificationService struct {
	repo  *repository.Repository
	cache *cache.Redis
}

func New(r *repository.Repository, c *cache.Redis) *NotificationService {
	return &NotificationService{
		repo:  r,
		cache: c,
	}
}
