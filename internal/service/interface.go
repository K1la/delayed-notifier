package service

import (
	"context"

	"github.com/K1la/delayed-notifier/internal/models"
)

type RepositoryI interface {
	CreateNotification(*models.Notification) (*models.Notification, error)
	GetNotifications(ctx context.Context) ([]models.Notification, error)
	GetNotificationStatusByID(ctx context.Context, id string) (string, error)
	CancelNotification(ctx context.Context, id string, newStatus models.NotificationStatus) error
}

type CacheI interface {
	Get(string) (string, error)
	Set(string, interface{}) error
}
