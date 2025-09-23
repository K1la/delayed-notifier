package main

import (
	"context"
	"github.com/K1la/delayed-notifier/internal/api/handlers"
	"github.com/K1la/delayed-notifier/internal/api/router"
	"github.com/K1la/delayed-notifier/internal/api/server"
	"github.com/K1la/delayed-notifier/internal/app"
	"github.com/K1la/delayed-notifier/internal/storage"
	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/zlog"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// context to handle SIGINT and SIGTERM for graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	zlog.Init()
	zlog.Logger.Info().Msgf("Start MAIN")

	val := validator.New()
	stor := storage.New()
	handler := handlers.New(stor, val)
	r := router.New(handler)
	s := server.New(":7777", r)

	worker := app.NewWorker(stor, 5*time.Second)
	worker.Start()
	defer worker.Stop()

	go func() {
		if err := s.ListenAndServe(); err != nil {
			zlog.Logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	<-ctx.Done()
	zlog.Logger.Info().Msg("shutdown signal received")
}
