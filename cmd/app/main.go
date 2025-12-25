package main

import (
	"log/slog"
	"os"

	"github.com/Telegram-bot-for-register-on-events/event-service/internal/config"
)

func main() {
	log := setupLogger()
	cfg := config.MustLoadConfig(log)
}

// TODO: Написать микросервис: сервер, принимающий запросы и возвращающий ответ

func setupLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
