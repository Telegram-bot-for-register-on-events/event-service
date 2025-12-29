package app

import (
	"log/slog"

	"github.com/Telegram-bot-for-register-on-events/event-service/internal/app/grpc"
	"github.com/Telegram-bot-for-register-on-events/event-service/internal/config"
	"github.com/Telegram-bot-for-register-on-events/event-service/internal/service"
	"github.com/Telegram-bot-for-register-on-events/event-service/internal/storage/postgres"
)

// App описывает микросервис целиком
type App struct {
	GRPCServer *grpcserver.App
	cfg        *config.Config
}

// NewApp конструктор для App
func NewApp(log *slog.Logger) *App {
	cfg := config.MustLoadConfig(log)
	db, err := postgres.NewStorage(log, cfg.GetDatabaseDriverName(), cfg.GetDatabasePath())
	if err != nil {
		panic(err)
	}
	s := service.NewService(log, db, db)
	grpcApp := grpcserver.New(log, cfg.GetGRPCServerPort(), s, s)
	return &App{
		GRPCServer: grpcApp,
		cfg:        cfg,
	}
}
