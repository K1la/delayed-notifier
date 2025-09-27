package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/wb-go/wbf/zlog"
)

func (s *NotificationService) handleMessage(ctx context.Context, message []byte, notification models.Notification) error {
	if err := json.Unmarshal(message, &notification); err != nil {
		return fmt.Errorf("failed to unmarshal notification from queue: %w", err)
	}

	// Пропускаем уведомления с failed или canceled статусом
	if notification.Status == models.StatusFailed || notification.Status == models.StatusCanceled || notification.Status == models.StatusSent {
		zlog.Logger.Info().Msgf("Skipping notification %s with status %s", notification.ID, notification.Status)
		return nil
	}

	// Если To содержит username (начинается с @), используем его как есть
	// Иначе пытаемся конвертировать в int
	var telegramId int
	var err error

	if len(notification.To) > 0 && notification.To[0] == '@' {
		// Для username пока что используем 0 (не будет работать с Telegram API)
		// В будущем нужно будет получать chatID по username
		telegramId = 0
		zlog.Logger.Warn().Msgf("Username %s cannot be used directly with Telegram API", notification.To)
		err = fmt.Errorf("username %s cannot be used directly", notification.To)
	} else {
		telegramId, err = strconv.Atoi(notification.To)
		if err != nil {
			zlog.Logger.Error().Err(err).Msgf("Failed to convert telegram id to int: %s", notification.To)
		}
	}

	// Если есть ошибка с telegramId, обрабатываем как неудачную попытку
	if err != nil {
		return s.handleFailedAttempt(ctx, notification, err)
	}

	// Пытаемся отправить сообщение
	if err := s.sender.SendToTelegram(telegramId, notification.Message); err != nil {
		zlog.Logger.Error().Err(err).Msgf("Failed to send notification to telegram: %s", notification.ID)
		return s.handleFailedAttempt(ctx, notification, err)
	}

	// Успешная отправка
	if err := s.repo.UpdateNotification(ctx, notification.ID, models.StatusSent); err != nil {
		return err
	}

	if err := s.cache.Set(notification.ID, string(models.StatusSent)); err != nil {
		return fmt.Errorf("failed to set in cache notification status to %s: %w", notification.ID, err)
	}

	zlog.Logger.Info().Msgf("Notification sent and stored in cache: %s", notification.ID)
	return nil
}

func (s *NotificationService) handleFailedAttempt(ctx context.Context, notification models.Notification, err error) error {
	// Увеличиваем счетчик попыток
	newRetries := notification.Retries + 1

	// Если попыток больше или равно 3, помечаем как failed
	if newRetries >= 3 {
		zlog.Logger.Error().Msgf("Notification %s failed after %d attempts, marking as failed", notification.ID, newRetries)

		if err := s.repo.UpdateNotification(ctx, notification.ID, models.StatusFailed); err != nil {
			return fmt.Errorf("failed to mark notification as failed: %w", err)
		}

		if err := s.cache.Set(notification.ID, string(models.StatusFailed)); err != nil {
			return fmt.Errorf("failed to set failed status in cache: %w", err)
		}

		return nil // Не возвращаем ошибку, чтобы сообщение не попало обратно в очередь
	}

	// Обновляем счетчик попыток в БД
	if err := s.repo.UpdateRetries(ctx, notification.ID, newRetries); err != nil {
		zlog.Logger.Error().Err(err).Msgf("Failed to update retries for notification %s", notification.ID)
	}

	zlog.Logger.Warn().Msgf("Notification %s failed (attempt %d/3), will retry later", notification.ID, newRetries)
	return err // Возвращаем ошибку, чтобы сообщение попало обратно в очередь для повторной попытки
}
