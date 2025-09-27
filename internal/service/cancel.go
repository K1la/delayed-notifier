package service

import (
	"context"
	"github.com/K1la/delayed-notifier/internal/models"
)

func (s *NotificationService) UpdateNotification(ctx context.Context, id string) error {
	if err := s.cache.Set(id, models.StatusCanceled); err != nil {
		return err
	}

	err := s.repo.UpdateNotification(ctx, id, models.StatusCanceled)
	if err != nil {
		return err
	}
	return nil
}
