package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/K1la/delayed-notifier/internal/api/handlers"
	"github.com/K1la/delayed-notifier/internal/api/router"
	"github.com/K1la/delayed-notifier/internal/api/server"
	"github.com/K1la/delayed-notifier/internal/cache"
	"github.com/K1la/delayed-notifier/internal/config"
	"github.com/K1la/delayed-notifier/internal/rabbitmq"
	"github.com/K1la/delayed-notifier/internal/repository"
	"github.com/K1la/delayed-notifier/internal/sender"
	"github.com/K1la/delayed-notifier/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/zlog"
)

func main() {

	zlog.Init()
	zlog.Logger.Info().Msgf("Start MAIN")

	cfg := config.Init()

	val := validator.New()
	db := repository.NewDB(cfg)
	repo := repository.New(db)

	ch := cache.New(cfg.Redis)
	queue := rabbitmq.New(cfg.RabbitMq)
	sender := sender.New()
	srvc := service.New(repo, ch, queue, sender)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// sig channel to handle SIGINT and SIGTERM for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Str("signal", sig.String()).Msg("Received shutdown signal. Shutting down...")
		cancel()
	}()

	go func() {
		if err := srvc.PublishPendingNotification(ctx); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("Failed to publish pending notification")
		}
	}()

	go func() {
		if err := srvc.ConsumeMessage(ctx); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("Failed to consume message")
		}
	}()

	handler := handlers.New(srvc, val)
	r := router.New(handler)
	s := server.New(cfg.HTTPServer.Address, r)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to start server")
		}
		zlog.Logger.Info().Msg("successfully started server on " + cfg.HTTPServer.Address)
	}()

	// Блокируем main горутину до получения сигнала завершения
	<-ctx.Done()
	zlog.Logger.Info().Msg("Shutting down gracefully...")
}
