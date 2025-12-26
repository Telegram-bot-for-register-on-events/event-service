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
	opLoadConfig = "config.load"
)

// Config описывает конфигурацию микросервиса
type Config struct {
	gRPCServerConfig *gRPCServerConfig
	databaseConfig   *databaseConfig
}

// gRPCServerConfig описывает конфигурацию gRPC-сервера
type gRPCServerConfig struct {
	port    string
	timeout time.Duration
}

// databaseConfig описывает конфигурацию базы данных
type databaseConfig struct {
	path string
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
		log.Error("error parsing GRPC_TIMEOUT", "timeout", timeoutStr, "error", err)
		return nil, err
	}

	if port == "" {
		log.Error("gRPC port cannot be empty")
		return nil, errors.New("gRPC port cannot be empty")
	}
	return &gRPCServerConfig{port, timeout}, nil
}

func newDatabaseConfig(log *slog.Logger) (*databaseConfig, error) {
	path := getEnv("DSN", "")
	if path == "" {
		log.Error("dsn cannot be empty")
		return nil, errors.New("dsn cannot be empty")
	}
	return &databaseConfig{path: path}, nil
}

// LoadConfig создаёт конфигурацию микросервиса
func LoadConfig(log *slog.Logger) (*Config, error) {
	log.Info("loading environment variables")
	// Загрузка переменных окружения из .env
	if err := godotenv.Load(); err != nil {
		log.Error("operation", opLoadConfig, err.Error())
		return nil, fmt.Errorf("%s: %w", opLoadConfig, err)
	}
	log.Info("environment variables successfully loaded")

	// Создаём конфигурацию базы данных
	dbCfg, err := newDatabaseConfig(log)
	if err != nil {
		log.Error("operation", opLoadConfig, err.Error())
		return nil, err
	}

	// Создаём конфигурацию gRPC-клиента
	gRPCCfg, err := newGRPCServerConfig(log)
	if err != nil {
		log.Error("operation", opLoadConfig, err.Error())
		return nil, err
	}

	return &Config{gRPCCfg, dbCfg}, nil
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
