package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/K1la/delayed-notifier/internal/models"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrNoNotificationsFound = errors.New("no notifications found")
)

//func (r *Repository) GetNotificationStatusById(ctx context.Context, id uuid.UUID) (string, error) {
//
//	return notif, nil
//}

func (r *Repository) GetNotifications(ctx context.Context) ([]models.Notification, error) {
	query := `
		SELECT id, message, send_at, retries, "to", channel, status
		FROM notifications
		ORDER BY send_at DESC;
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		if err = rows.Scan(&n.ID, &n.Message, &n.SendAt, &n.Retries, &n.Channel, &n.Status); err != nil {
			return nil, err
		}

		notifications = append(notifications, n)
	}

	if len(notifications) == 0 {
		return nil, ErrNoNotificationsFound
	}

	return notifications, nil
}
