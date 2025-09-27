package handlers

import (
	"context"

	"github.com/K1la/delayed-notifier/internal/models"
)

type ServiceI interface {
	CreateNotification(context.Context, *models.Notification) (*models.Notification, error)
	GetNotifications(context.Context) ([]models.Notification, error)
	GetNotificationStatusByID(context.Context, string) (string, error)
	UpdateNotification(context.Context, string) error
}
