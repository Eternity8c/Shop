package app

import (
	grpcapp "auth-service/internal/app/grpc"
	"auth-service/internal/services/auth"
	"auth-service/internal/storage/postgres"
	"context"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, stroagePath string, tokenTTL time.Duration, secret string) *App {
	storage, err := postgres.New(context.Background(), stroagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, tokenTTL, secret)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
