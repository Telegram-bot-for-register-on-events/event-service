package app

import (
	"log/slog"
	"os"

	"github.com/Telegram-bot-for-register-on-events/event-service/internal/app/grpc"
	"github.com/Telegram-bot-for-register-on-events/event-service/internal/config"
	"github.com/Telegram-bot-for-register-on-events/event-service/internal/nats"
	"github.com/Telegram-bot-for-register-on-events/event-service/internal/service"
	"github.com/Telegram-bot-for-register-on-events/event-service/internal/storage/postgres"
)

// App описывает микросервис целиком
type App struct {
	log        *slog.Logger
	GRPCServer *grpcserver.App
	cfg        *config.Config
	Nats       *nats.Nats
	Database   *postgres.Storage
}

// NewApp конструктор для App
func NewApp(log *slog.Logger) *App {
	// Загружаем конфигурацию
	cfg := cfgInit(log)
	// Инициализируем хранилище данных
	db := dbInit(log, cfg.GetDatabaseDriverName(), cfg.GetDatabasePath())
	// Инициализируем сервисный слой
	s := service.NewService(log, db, db)
	// Подключаемся к Nats
	n := natsConn(log, cfg.GetNatsURL())
	// Создаём поток и топик
	_, err := n.CreateStream(cfg.GetNatsStream(), []string{cfg.GetNatsTopic()})
	if err != nil {
		log.Error("error", err.Error(), slog.String("failed", "create stream in NATS"))
		os.Exit(1)
	}
	// Создаём gRPC-сервер
	grpcApp := grpcserver.New(log, cfg.GetGRPCServerPort(), s, n, s)

	return &App{
		log:        log,
		GRPCServer: grpcApp,
		cfg:        cfg,
		Nats:       n,
		Database:   db,
	}
}

// MustStart запускает микросервис
func (a *App) MustStart() {
	a.log.Info("application successfully started")
	go a.GRPCServer.MustRun()
}

// Stop выполняет остановку всего микросервиса
func (a *App) Stop() {
	a.log.Info("shutting down...")
	a.GRPCServer.Stop()
	a.Nats.Conn.Close()
	a.Database.Close()
}

// cfgInit обёртка для инициализации конфига
func cfgInit(log *slog.Logger) *config.Config {
	cfg := config.MustLoadConfig(log)
	log.Info("config successfully loaded")
	return cfg
}

// dbInit обёртка для создания подключения к базе данных
func dbInit(log *slog.Logger, driverName, dsn string) *postgres.Storage {
	db, err := postgres.NewStorage(log, driverName, dsn)
	if err != nil {
		log.Error("error", err.Error(), slog.String("failed", "connect to database"))
		os.Exit(1)
	}
	log.Info("connection to database successfully")
	return db
}

// natsConn обёртка для подключения к NATS
func natsConn(log *slog.Logger, url string) *nats.Nats {
	n, err := nats.NewNats(log, url)
	if err != nil {
		log.Error("error", err.Error(), slog.String("failed", "connect to nats"))
		os.Exit(1)
	}
	log.Info("connection to nats successfully")
	return n
}
