package handlers

import (
	"context"

	"github.com/K1la/delayed-notifier/internal/models"
)

type ServiceI interface {
	CreateNotification(ctx context.Context, n *models.Notification) (*models.Notification, error)
	GetNotifications(ctx context.Context) ([]models.Notification, error)
	GetNotificationStatusByID(ctx context.Context, id string) (string, error)
	CancelNotification(ctx context.Context, id string) error
}
