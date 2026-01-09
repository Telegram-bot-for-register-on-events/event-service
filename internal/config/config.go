package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Константы для описания операций
const (
	opLoadConfig      = "config.load"
	opNewServerConfig = "config.NewGRPCServerConfig"
)

// Config описывает конфигурацию микросервиса
type Config struct {
	gRPCServerConfig *gRPCServerConfig
	databaseConfig   *databaseConfig
	natsConfig       *natsConfig
}

// gRPCServerConfig описывает конфигурацию gRPC-сервера
type gRPCServerConfig struct {
	port    string
	timeout time.Duration
}

// databaseConfig описывает конфигурацию базы данных
type databaseConfig struct {
	driverName string
	path       string
}

// natsConfig описывает конфигурацию NATS
type natsConfig struct {
	url    string
	stream string
	topic  string
}

// getEnv проверяет наличие переменной окружения и возвращает её текущее значение, либо стандартное, при отсутствии текущего
func getEnv(key, reserve string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return reserve
}

// newGRPCServerConfig загружает конфигурацию для gRPC-сервера
func newGRPCServerConfig(log *slog.Logger) (*gRPCServerConfig, error) {
	port := getEnv("GRPC_PORT", "")
	timeoutStr := getEnv("GRPC_TIMEOUT", "")
	if timeoutStr == "" {
		log.Error("gRPC timeout cannot be empty")
		return nil, errors.New("gRPC timeout cannot be empty")
	}

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opNewServerConfig))
		return nil, err
	}

	if port == "" {
		log.Error("gRPC port cannot be empty")
		return nil, errors.New("gRPC port cannot be empty")
	}
	return &gRPCServerConfig{port, timeout}, nil
}

// newDatabaseConfig загружает конфигурацию для базы данных
func newDatabaseConfig(log *slog.Logger) (*databaseConfig, error) {
	path := getEnv("DSN", "")
	if path == "" {
		log.Error("dsn cannot be empty")
		return nil, errors.New("dsn cannot be empty")
	}
	driverName := getEnv("DB_DRIVER_NAME", "")
	if driverName == "" {
		log.Error("database driver name cannot be empty")
		return nil, errors.New("database driver name cannot be empty")
	}
	return &databaseConfig{driverName: driverName, path: path}, nil
}

// newNatsConfig загружает конфигурацию для NATS
func newNatsConfig(log *slog.Logger) (*natsConfig, error) {
	url := getEnv("NATS_URL", "")
	if url == "" {
		log.Error("nats url cannot be empty")
		return nil, errors.New("nats url cannot be empty")
	}

	topic := getEnv("NATS_TOPIC", "")
	if topic == "" {
		log.Error("nats topic cannot be empty")
		return nil, errors.New("nats topic cannot be empty")
	}

	stream := getEnv("NATS_STREAM", "")
	if stream == "" {
		log.Error("nats stream cannot be empty")
		return nil, errors.New("nats stream cannot be empty")
	}
	return &natsConfig{url: url, stream: stream, topic: topic}, nil
}

// LoadConfig создаёт конфигурацию микросервиса
func LoadConfig(log *slog.Logger) (*Config, error) {
	log.Info("loading environment variables")
	// Загрузка переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Error("error", err.Error(), slog.String("operation", opLoadConfig))
		return nil, fmt.Errorf("%s: %w", opLoadConfig, err)
	}
	log.Info("environment variables successfully loaded")

	// Создаём конфигурацию базы данных
	dbCfg, err := newDatabaseConfig(log)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opLoadConfig))
		return nil, fmt.Errorf("%s: %w", opLoadConfig, err)
	}

	// Создаём конфигурацию gRPC-клиента
	gRPCCfg, err := newGRPCServerConfig(log)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opLoadConfig))
		return nil, fmt.Errorf("%s: %w", opLoadConfig, err)
	}

	natsCfg, err := newNatsConfig(log)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opLoadConfig))
		return nil, fmt.Errorf("%s: %w", opLoadConfig, err)
	}

	return &Config{gRPCCfg, dbCfg, natsCfg}, nil
}

// MustLoadConfig обёртка для LoadConfig, при ошибке - паникует
func MustLoadConfig(log *slog.Logger) *Config {
	cfg, err := LoadConfig(log)
	if err != nil {
		panic(err)
	}
	return cfg
}

// GetGRPCServerPort геттер, для получения порта gRPC-сервера
func (c *Config) GetGRPCServerPort() string {
	return c.gRPCServerConfig.port
}

// GetDatabasePath геттер, для получения пути подключения к базе данных
func (c *Config) GetDatabasePath() string {
	return c.databaseConfig.path
}

// GetDatabaseDriverName геттер для получения драйвера базы данных
func (c *Config) GetDatabaseDriverName() string {
	return c.databaseConfig.driverName
}

// GetNatsURL геттер для получения URL для подключения к NATS
func (c *Config) GetNatsURL() string {
	return c.natsConfig.url
}

// GetNatsStream геттер для получения названия для создания потока в NATS
func (c *Config) GetNatsStream() string {
	return c.natsConfig.stream
}

// GetNatsTopic геттер для получения названия топика, в который будут публиковаться сообщения
func (c *Config) GetNatsTopic() string { return c.natsConfig.topic }
