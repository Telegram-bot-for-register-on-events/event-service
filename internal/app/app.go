package app

import (
	"log/slog"

	"github.com/Telegram-bot-for-register-on-events/event-service/internal/app/grpc"
	"github.com/Telegram-bot-for-register-on-events/event-service/internal/config"
)

type App struct {
	GRPCServer *grpcserver.App
	cfg        *config.Config
}

func NewApp(log *slog.Logger) *App {
	cfg := config.MustLoadConfig(log)
	grpcApp := grpcserver.New(log, cfg.GetGRPCServerPort())
	return &App{
		GRPCServer: grpcApp,
		cfg:        cfg,
	}
}
