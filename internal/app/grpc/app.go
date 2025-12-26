package app

import (
	"log/slog"

	eventgrpc "github.com/Telegram-bot-for-register-on-events/event-service/internal/grpc"
	"google.golang.org/grpc"
)

// App описывает серверное gRPC-приложение
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

// NewApp создаёт новое серверное gRPC-приложение
func NewApp(log *slog.Logger, port string) *App {
	// Создаём новый gRPC-сервер
	grpcServer := grpc.NewServer()
	// Подключаем к нему обработчик
	eventgrpc.Register(grpcServer)

	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}
