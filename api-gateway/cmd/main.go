package main

import (
	"api-geteway/internal/auth"
	"api-geteway/internal/config"
	authHandlers "api-geteway/internal/handlers/auth"
	"context"
	"log/slog"
	"net/http"
	"os"
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

	mux := http.NewServeMux()
	ah := authHandlers.NewAuthHandler(client, logger)

	mux.HandleFunc("/login", ah.Login)
	mux.HandleFunc("/register", ah.Register)

	logger.Info("starting api-gateway", "addr", ":"+cfg.APIGatewayPort)

	err = http.ListenAndServe(":"+cfg.APIGatewayPort, mux)
	if err != nil {
		logger.Warn("listen and serve:", "err", err)
	}
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
