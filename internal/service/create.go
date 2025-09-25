package service

import (
	"context"

	"github.com/K1la/delayed-notifier/internal/models"
)

func (s *NotificationService) CreateNotification(ctx context.Context, n *models.Notification) (*models.Notification, error) {
	notif, err := s.repo.CreateNotification(ctx, n)
	if err != nil {
		return nil, err
	}

	if err = s.cache.Set(notif.ID, notif.Status); err != nil {
		return nil, err
	}
	return notif, nil
}
