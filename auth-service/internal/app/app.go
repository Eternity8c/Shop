package app

import (
	grpcapp "auth-service/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPC *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, stroagePath string, tokenTTL time.Duration) *App {
	// TODO: доделать

	panic("implement me")
}
