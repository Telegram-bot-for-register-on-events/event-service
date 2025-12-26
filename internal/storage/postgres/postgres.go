package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Константы для описания операций
const (
	opConnect         = "postgres.connect"
	opCloseConnection = "postgres.closeConnection"
)

// Storage описывает слой взаимодействия с базой данных
type Storage struct {
	DB  *sqlx.DB
	log *slog.Logger
}

// NewStorage описывает объект базы данных
func NewStorage(log *slog.Logger, driverName, dsn string) (*Storage, error) {
	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		log.Error("operation", opConnect, err.Error())
		return nil, fmt.Errorf("%s: %w", opConnect, err)
	}

	// Проверяем подключение к базе данных, в противном случае возвращаем ошибку
	if err = db.Ping(); err != nil {
		log.Error("operation", opConnect, err.Error())
		return nil, fmt.Errorf("%s: %w", opConnect, err)
	}
	log.Info("connection to database successfully")
	return &Storage{
		DB:  db,
		log: log,
	}, nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() {
	s.log.Info("operation", opCloseConnection)
	if err := s.DB.Close(); err != nil {
		s.log.Error("closing database connection", err.Error())
	}
}

func (s *Storage) GetEvents(ctx context.Context) ([]*event.Event, error) {
	panic("storage.GetEvents not implemented")
}

func (s *Storage) GetEvent(ctx context.Context, eventID string) (*event.Event, error) {
	panic("storage.GetEvent not implemented")
}
