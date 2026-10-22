package main

import (
	"api-geteway/internal/client/auth"
	client_cart "api-geteway/internal/client/cart"
	"api-geteway/internal/client/product"
	"api-geteway/internal/config"
	authHandlers "api-geteway/internal/handlers/auth"
	handler_cart "api-geteway/internal/handlers/cart"
	productHandler "api-geteway/internal/handlers/product"
	MyMetrics "api-geteway/internal/metrics"
	"api-geteway/internal/router"
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg := config.MustLoad()
	logger := setupLogger()

	addrAuth := cfg.AuthServicesAddr

	client, err := auth.New(context.Background(), logger, addrAuth)
	if err != nil {
		logger.Warn("create auth client:", "err", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			logger.Warn("client.Close:", "err", err)
		}
	}()

	productClient := product.New(cfg.ProductServiceAddr, logger)

	cartClient := client_cart.New(cfg.CartServiceAddr, logger)

	mux := http.NewServeMux()

	ah := authHandlers.NewAuthHandler(client, logger)
	ph := productHandler.New(productClient, logger)
	clienHandler := handler_cart.NewCartHandler(*cartClient, logger)

	router := router.New(ah, ph, clienHandler)

	router.RegisterRoutes(mux)

	mux.Handle("/metrics", promhttp.Handler())

	handler := MyMetrics.InstrumentHandler(mux)

	srv := &http.Server{
		Addr:    ":" + cfg.APIGatewayPort,
		Handler: handler,
	}

	logger.Info("starting api-gateway", "addr", ":"+cfg.APIGatewayPort)

	err = srv.ListenAndServe()
	if err != nil {
		logger.Warn("listen and serve:", "err", err)
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
