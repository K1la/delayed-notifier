package service

import (
	"context"
	"github.com/K1la/delayed-notifier/internal/models"
)

func (s *NotificationService) CancelNotification(ctx context.Context, id string) error {
	if err := s.cache.Set(id, models.StatusCanceled); err != nil {
		return err
	}

	err := s.repo.CancelNotification(ctx, id, models.StatusCanceled)
	if err != nil {
		return err
	}
	return nil
}
