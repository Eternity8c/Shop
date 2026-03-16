package main

import (
	"cart-service/internal/config"
	"cart-service/internal/storage"
	"context"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	logger := setupLogger()

	_, err := storage.NewStorage(ctx, cfg.StoragePath, logger)
	if err != nil {
		logger.Error("Storage error", "error", err)
		return
	}

	logger.Info("service is start", "Port:", cfg.Port)
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
