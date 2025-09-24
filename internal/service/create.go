package service

import (
	"context"
	"github.com/K1la/delayed-notifier/internal/models"
)

func (s *NotificationService) CreateNotification(ctx context.Context, n *models.Notification) error {
	err := s.repo.CreateNotification(ctx, n)
	if err != nil {
		return err
	}

	if err = s.cache.Set(n.ID, n.Status); err != nil {
		return err
	}
	return nil
}
