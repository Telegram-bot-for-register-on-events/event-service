package grpcserver

import (
	"fmt"
	"log/slog"
	"net"

	eventgrpc "github.com/Telegram-bot-for-register-on-events/event-service/internal/grpc/event"
	"google.golang.org/grpc"
)

// Константы для описания операций
const (
	opStart = "grpcserver.Start"
	opStop  = "grpcserver.Stop"
)

// App описывает gRPC-сервер
type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       string
}

// New создаёт новый gRPC-сервер
func New(log *slog.Logger, port string, events eventgrpc.EventService) *App {
	grpcServer := grpc.NewServer()
	// Подключаем обработчик
	eventgrpc.Register(grpcServer, events)
	return &App{
		log:        log,
		gRPCServer: grpcServer,
		port:       port,
	}
}

// start запускает gRPC-сервер
func (a *App) start() error {
	// Инициализируем tcp-слушателя на указанном порту
	listener, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return fmt.Errorf("error starting gRPC server - %w", err)
	}

	a.log.Info("running gRPC-server...", slog.String("operation", opStart), slog.String("port", a.port))

	// Принимаем входящие соединения от слушателя
	if err = a.gRPCServer.Serve(listener); err != nil {
		a.log.Info("error starting gRPC server", slog.String("error", err.Error()))
		return fmt.Errorf("error starting gRPC server - %w", err)
	}

	return nil
}

// MustRun обёртка для метода start, при ошибке - паникует
func (a *App) MustRun() {
	if err := a.start(); err != nil {
		panic(err)
	}
}

// Stop останавливает gRPC-сервер
func (a *App) Stop() {
	a.log.Info("gRPC-server is stopping", slog.String("operation", opStop), slog.String("port", a.port))
	a.gRPCServer.GracefulStop()
}
