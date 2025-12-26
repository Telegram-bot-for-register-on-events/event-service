package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Telegram-bot-for-register-on-events/event-service/internal/domain/models"
	"github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Константы для описания операций
const (
	opConnect         = "postgres.connect"
	opCloseConnection = "postgres.closeConnection"
	opGetEvents       = "postgres.getEvents"
	opGetEvent        = "postgres.getEvent"
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
	var eventsDB []models.Event
	err := s.DB.SelectContext(ctx, &eventsDB, `select * from events`)
	if err != nil {
		s.log.Error("operation", opGetEvents, err.Error())
		return nil, fmt.Errorf("%s: %w", opGetEvents, err)
	}
	// Преобразуем DB-структуры в protobuf-структуры
	events := make([]*event.Event, 0, len(eventsDB))
	for _, e := range eventsDB {
		events = append(events, convertingEventsStruct(e))
	}
	return events, nil
}

func (s *Storage) GetEvent(ctx context.Context, eventID string) (*event.Event, error) {
	var e models.Event
	err := s.DB.GetContext(ctx, &e, `select * from events where id = $1`, eventID)
	if err != nil {
		s.log.Error("operation", opGetEvent, err.Error())
		return nil, fmt.Errorf("%s: %w", opGetEvent, err)
	}
	return convertingEventsStruct(e), nil
}

func convertingEventsStruct(eventDB models.Event) *event.Event {
	return &event.Event{
		Id:          eventDB.ID,
		Title:       eventDB.Title,
		Description: eventDB.Description,
		StartsAt:    timestamppb.New(eventDB.StartsAt),
	}
}
