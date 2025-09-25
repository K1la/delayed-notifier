package service

import (
	"context"

	"github.com/K1la/delayed-notifier/internal/models"
)

type Storage interface {
	CreateNotification(*models.Notification) (*models.Notification, error)
	GetNotifications(ctx context.Context) ([]models.Notification, error)
}

type Cache interface {
	Get(string) (string, error)
	Set(string, interface{}) error
}
