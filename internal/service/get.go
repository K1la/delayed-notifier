package service

import (
	"context"
	"errors"

	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/wb-go/wbf/zlog"
)

func (s *NotificationService) GetNotifications(ctx context.Context) ([]models.Notification, error) {
	return s.repo.GetNotifications(ctx)
}

func (s *NotificationService) GetNotificationStatusByID(ctx context.Context, id string) (string, error) {
	status, err := s.cache.Get(id)
	if err != nil && !errors.Is(err, redis.Nil) {
		zlog.Logger.Error().Err(err).Str("id", id).Msg("failed to get notification status from cache")
	}

	if errors.Is(err, redis.Nil) {
		status, err = s.repo.GetNotificationStatusById(ctx, id)
		if err != nil {
			return "", err
		}

		err = s.cache.Set(id, status)
		if err != nil {
			zlog.Logger.Error().Err(err).Str("id", id).Msg("failed to cache notification status")
		}
	}

	return status, nil
}
