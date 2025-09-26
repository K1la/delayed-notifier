package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/K1la/delayed-notifier/internal/models"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrNoNotificationsFound = errors.New("no notifications found")
)

func (r *Repository) GetNotificationStatusByID(ctx context.Context, id string) (string, error) {
	query := `
		SELECT status
		FROM notifications
		WHERE id = $1
	`

	var status string
	err := r.db.Master.QueryRowContext(ctx, query, id).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotificationNotFound
		}
		return "", fmt.Errorf("failed to get notification status from db: %w", err)
	}

	return status, nil
}

func (r *Repository) GetNotifications(ctx context.Context) ([]models.Notification, error) {
	query := `
		SELECT id, message, send_at, retries, "to", channel, status, created_at, updated_at
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
		err = rows.Scan(
			&n.ID,
			&n.Message,
			&n.SendAt,
			&n.Retries,
			&n.To,
			&n.Channel,
			&n.Status,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row to model: %w", err)
		}

		notifications = append(notifications, n)
	}

	if len(notifications) == 0 {
		return nil, ErrNoNotificationsFound
	}

	return notifications, nil
}

func (r *Repository) GetPendingNotifications(ctx context.Context) ([]models.Notification, error) {
	query := `
	SELECT id, message, send_at, retries, "to", channel, status, created_at, updated_at
	FROM notifications
	WHERE send_at < NOW()
	AND status = 'pending'
	`

	rows, err := r.db.Master.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all pending notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err = rows.Scan(
			&n.ID,
			&n.Message,
			&n.SendAt,
			&n.Retries,
			&n.To,
			&n.Channel,
			&n.Status,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row to model: %w", err)
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}
