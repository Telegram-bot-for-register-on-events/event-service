package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Telegram-bot-for-register-on-events/event-service/internal/app"
)

func main() {
	// Инициализируем логгер
	log := setupLogger()
	// Создаём новый инстанс микросервиса
	application := app.NewApp(log)
	// Запускаем его
	go application.GRPCServer.MustRun()
	// Создаём канал для приёма сигналов операционной системы
	stop := make(chan os.Signal, 1)
	// Передаём входящие сигналы в канал stop
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	// Читаем из канала, пока не придёт соответствующий сигнал
	<-stop

	application.GRPCServer.Stop()
}

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
