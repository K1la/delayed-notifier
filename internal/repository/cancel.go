package repository

import (
	"context"
	"fmt"
	"github.com/K1la/delayed-notifier/internal/models"
)

func (r *Repository) CancelNotification(ctx context.Context, id string, newStatus models.NotificationStatus) error {
	query := `
	UPDATE notification
	SET status = $1
	WHRE id = $2
	`

	res, err := r.db.Master.ExecContext(ctx, query, newStatus, id)
	if err != nil {
		return fmt.Errorf("failed to update notification status to %s: %w", newStatus, err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get number of affected rows: %w", err)
	}

	if affected == 0 {
		return ErrNotificationNotFound
	}

	return nil
}
