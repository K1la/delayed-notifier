package repository

import (
	"context"
	"fmt"
	"github.com/wb-go/wbf/zlog"

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
		zlog.Logger.Error().Err(err).Msgf("Failed to create notification id=%d", notif.ID)
		return nil, fmt.Errorf("could not scan notification from db: %w", err)
	}
	zlog.Logger.Debug().Msgf("Created notification id=%d", notif.ID)
	return notif, nil
}
