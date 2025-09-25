package repository

import (
	"context"
	"fmt"
	"github.com/K1la/delayed-notifier/internal/models"
)

func (r *Repository) CreateNotification(ctx context.Context, notif *models.Notification) (*models.Notification, error) {
	query := `
	INSERT INTO notifications (
		message, send_at, retries, "to", channel
	) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
	`

	err := r.db.Master.QueryRowContext(
		ctx,
		query,
		notif.Message,
		notif.SendAt,
		notif.Retries,
		notif.To,
		notif.Channel,
	).Scan(&notif.ID, &notif.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not scan notification from db: %w", err)
	}

	return notif, nil
}
