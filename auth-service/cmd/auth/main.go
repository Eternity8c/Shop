package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Info("starting application", slog.Any("config", cfg))

	application := app.New(log, cfg.Port, cfg.StoragePath, cfg.TokenTTL, cfg.Secret)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}
