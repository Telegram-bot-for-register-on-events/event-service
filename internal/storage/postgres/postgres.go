package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
	opRegister        = "postgres.register"
)

// Storage описывает слой взаимодействия с базой данных
type Storage struct {
	DB  *sqlx.DB
	log *slog.Logger
}

// NewStorage конструктор для Storage
func NewStorage(log *slog.Logger, driverName, dsn string) (*Storage, error) {
	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		log.Error("error", err.Error(), slog.String("operation", opConnect))
		return nil, fmt.Errorf("%s: %w", opConnect, err)
	}

	// Проверяем подключение к базе данных, в противном случае возвращаем ошибку
	if err = db.Ping(); err != nil {
		log.Error("error", err.Error(), slog.String("operation", opConnect))
		return nil, fmt.Errorf("%s: %w", opConnect, err)
	}

	return &Storage{
		DB:  db,
		log: log,
	}, nil
}

// Close закрывает соединение с базой данных
func (s *Storage) Close() {
	if err := s.DB.Close(); err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opCloseConnection))
	}
}

func (s *Storage) GetEvents(ctx context.Context) ([]*event.Event, error) {
	var eventsDB []models.Event
	err := s.DB.SelectContext(ctx, &eventsDB, `select * from events`)
	if err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opGetEvents))
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
		s.log.Error("error", err.Error(), slog.String("operation", opGetEvent))
		return nil, fmt.Errorf("%s: %w", opGetEvent, err)
	}
	return convertingEventsStruct(e), nil
}

func (s *Storage) RegisterUser(ctx context.Context, eventID string, chatID int64, username string) error {
	query := `insert into registration (event_id, chat_id, username, created_at) values ($1, $2, $3, $4)`
	reg := &models.Registration{
		EventID:   eventID,
		ChatID:    chatID,
		Username:  username,
		CreatedAt: time.Now(),
	}
	_, err := s.DB.NamedExecContext(ctx, query, reg)
	if err != nil {
		s.log.Error("error", err.Error(), slog.String("operation", opRegister))
		return fmt.Errorf("%s: %w", opRegister, err)
	}
	return nil
}

func convertingEventsStruct(eventDB models.Event) *event.Event {
	return &event.Event{
		Id:          eventDB.ID,
		Title:       eventDB.Title,
		Description: eventDB.Description,
		StartsAt:    timestamppb.New(eventDB.StartsAt),
	}
}
