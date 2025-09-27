package service

import (
	"context"

	"github.com/K1la/delayed-notifier/internal/models"
)

type RepositoryI interface {
	CreateNotification(context.Context, *models.Notification) (*models.Notification, error)
	GetNotifications(context.Context) ([]models.Notification, error)
	GetNotificationStatusByID(context.Context, string) (string, error)
	GetPendingNotifications(context.Context) ([]models.Notification, error)
	UpdateNotification(context.Context, string, models.NotificationStatus) error
	UpdateRetries(context.Context, string, int) error
}

type CacheI interface {
	Get(string) (string, error)
	Set(string, interface{}) error
}

type QueueI interface {
	Publish(models.Notification) error
	Consume(context.Context) (<-chan []byte, error)
}

type SenderI interface {
	SendToTelegram(telegramId int, message string) error
}
