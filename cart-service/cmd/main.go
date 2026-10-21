package main

import (
	"cart-service/internal/client"
	"cart-service/internal/config"
	"cart-service/internal/handler"
	"cart-service/internal/router"
	"cart-service/internal/service"
	"cart-service/internal/storage"
	"context"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	logger := setupLogger()

	storage, err := storage.NewStorage(ctx, cfg.StoragePath, logger)
	if err != nil {
		logger.Error("Storage error", "error", err)
		return
	}
	client := client.New("", logger)
	service := service.NewCartService(logger, storage, client)
	handler := handler.NewCartHandler(logger, service)

	mux := http.NewServeMux()

	r := router.New(handler)
	r.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	logger.Info("service is start", "Port:", cfg.Port)

	if err := srv.ListenAndServe(); err != nil {
		logger.Error("listen and serve:", "err", err)
	}

}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
