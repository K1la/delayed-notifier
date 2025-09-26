package handlers

import (
	"context"

	"github.com/K1la/delayed-notifier/internal/models"
)

type ServiceI interface {
	GetNotifications(ctx context.Context) ([]*models.Notification, error)
	CreateNotification(ctx context.Context, n *models.Notification) error
	GetNotificationStatusByID(ctx context.Context, id string) (*models.NotificationStatus, error)
	CancelNotification(ctx context.Context, id string) error
}
