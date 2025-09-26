package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"time"
)

func (r *RabbitMQ) Publish(notification models.Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to marshal notification ot send to rabbitmq")
	}

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  2,
	}
	return r.publisher.PublishWithRetry(body, "notification", "application/json", strategy)
}

func (r *RabbitMQ) Consume(ctx context.Context) (<-chan []byte, error) {
	messagesCh := make(chan []byte)

	go func() {
		defer close(messagesCh)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				next, ok := <-r.consumer
				if !ok {
					return
				}

				if err := next.Ack(false); err != nil {
					zlog.Logger.Error().Err(err).Msg("failed to acknowledge message")
				}
				messagesCh <- next.Body
			}
		}
	}()

	return messagesCh, nil
}
