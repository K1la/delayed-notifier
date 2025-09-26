package rabbitmq

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"

	"github.com/K1la/delayed-notifier/internal/config"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/zlog"
)

type RabbitMQ struct {
	publisher *rabbitmq.Publisher
	consumer  <-chan amqp091.Delivery
}

func New(cfg config.RabbitMqCfg) *RabbitMQ {
	url := fmt.Sprintf(
		"amqp://guest:guest@%s:%s/",
		cfg.Host, cfg.Port,
	)

	connection, err := rabbitmq.Connect(url, cfg.Retries, cfg.Pause)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to connect to RabbitMQ")
	}

	pubCh, err := connection.Channel()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to create channel for publisher rabbitmq")
	}

	qm := rabbitmq.NewQueueManager(pubCh)
	_, err = qm.DeclareQueue("notification")
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to declare queue")
	}

	publisher := rabbitmq.NewPublisher(pubCh, "")

	conCh, err := connection.Channel()
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("Failed to create channel for consumer rabbitmq")
	}

	deliveries, err := conCh.Consume("notification", "", false, false, false, false, amqp091.Table{})
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to create consumer for rabbitmq")
	}

	return &RabbitMQ{publisher: publisher, consumer: deliveries}
}
