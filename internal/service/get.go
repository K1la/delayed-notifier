package service

import (
	"context"
	"fmt"
	"github.com/K1la/delayed-notifier/internal/models"
)

func (s *NotificationService) GetNotifications(ctx context.Context) ([]models.Notification, error) {
	notifications, err := s.repo.GetNotifications(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all notifications: %w", err)
	}

	return notifications, nil
}
