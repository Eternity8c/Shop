package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"product-service/internal/config"
	"product-service/internal/handler"
	"product-service/internal/router"
	"product-service/internal/service"
	"product-service/internal/storage"
)

func main() {
	cfg := config.MustLoad()
	logger := setupLogger()
	ctx := context.Background()

	store, err := storage.New(ctx, cfg.StoragePath)
	if err != nil {
		logger.Warn("failed to connect db", err)
	}
	svc := service.New(logger, store)
	h := handler.New(ctx, logger, svc)
	r := router.New(h)

	mux := http.NewServeMux()
	r.RegisterRoutes(mux)

	srv := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	logger.Info("server is starting", "Port:", cfg.Port)

	if err := srv.ListenAndServe(); err != nil {
		logger.Warn("err", err)
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
