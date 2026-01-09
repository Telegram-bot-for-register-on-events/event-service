package service

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
)

// Константы для описания операций
const (
	opGetEvents = "service.GetEvents"
	opGetEvent  = "service.GetEvent"
	opRegister  = "service.Register"
)

// Service описывает сервисный слой микросервиса
type Service struct {
	log           *slog.Logger
	eventReceiver EventReceiver
	registerer    Registerer
}

// EventReceiver описывает методы для получения информации о событиях
type EventReceiver interface {
	GetEvents(ctx context.Context) ([]*pb.Event, error)
	GetEvent(ctx context.Context, eventID string) (*pb.Event, error)
}

// Registerer описывает метод для взаимодействия с repo-слоем
type Registerer interface {
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) error
}

// NewService конструктор для создания Service
func NewService(log *slog.Logger, eventReceiver EventReceiver, registerer Registerer) *Service {
	return &Service{
		log:           log,
		eventReceiver: eventReceiver,
		registerer:    registerer,
	}
}

func (s *Service) GetEvents(ctx context.Context) ([]*pb.Event, error) {
	events, err := s.eventReceiver.GetEvents(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetEvents, err)
	}
	return events, nil
}

func (s *Service) GetEvent(ctx context.Context, eventID string) (*pb.Event, error) {
	event, err := s.eventReceiver.GetEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opGetEvent, err)
	}
	return event, nil
}

func (s *Service) RegisterUser(ctx context.Context, eventID string, chatID int64, username string) error {
	err := s.registerer.RegisterUser(ctx, eventID, chatID, username)
	if err != nil {
		return fmt.Errorf("%s: %w", opRegister, err)
	}
	return nil
}
