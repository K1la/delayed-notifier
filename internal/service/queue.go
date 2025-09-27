package service

import (
	"context"
	"time"

	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/wb-go/wbf/zlog"
)

const (
	workersN = 3
)

func (s *NotificationService) PublishPendingNotification(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second) // Проверяем каждые 30 секунд
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			notifications, err := s.repo.GetPendingNotifications(ctx)
			if err != nil {
				zlog.Logger.Error().Err(err).Msg("Failed to get pending notifications")
				continue
			}
			for _, n := range notifications {
				// Пропускаем уведомления с failed или canceled статусом
				if n.Status == models.StatusFailed || n.Status == models.StatusCanceled {
					zlog.Logger.Info().Msgf("Skipping notification %s with status %s", n.ID, n.Status)
					continue
				}

				// Пропускаем уведомления, которые еще не готовы к отправке
				if (n.SendAt.UnixMilli() - time.Now().UnixMilli()) > time.Minute.Milliseconds()/2 {
					continue
				}

				if err = s.queue.Publish(n); err != nil {
					zlog.Logger.Error().Err(err).Msg("Failed to publish notification to queue")
					continue
				}
				zlog.Logger.Debug().Msgf("Queue publish notification %+v", n)
			}
		}
	}
}

func (s *NotificationService) ConsumeMessage(ctx context.Context) error {
	messages, err := s.queue.Consume(ctx)
	if err != nil {
		return err
	}

	for i := range workersN {
		go func(i int) {
			zlog.Logger.Info().Msgf("Worker #%d start consuming", i)
			for msg := range messages {
				var n models.Notification
				if err = s.handleMessage(ctx, msg, n); err != nil {
					zlog.Logger.Error().Err(err)
					continue
				}
			}
		}(i)
	}
	return nil
}
