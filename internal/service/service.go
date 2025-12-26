package service

import (
	"context"
	"fmt"
	"log/slog"

	pb "github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
)

// Константы для описания операций
const (
	opGetEvents    = "service.GetEvents"
	opGetEvent     = "service.GetEvent"
	opRegisterUser = "service.RegisterUser"
)

// Service описывает сервисный слой микросервиса
type Service struct {
	log           *slog.Logger
	eventReceiver EventReceiver
	userRegister  UserRegister
}

// EventReceiver описывает методы для получения информации о событиях
type EventReceiver interface {
	GetEvents(ctx context.Context) ([]*pb.Event, error)
	GetEvent(ctx context.Context, eventID string) (*pb.Event, error)
}

// UserRegister описывает метод для регистрации пользователя на конкретное событие
type UserRegister interface {
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error)
}

// NewService конструктор для создания Service
func NewService(log *slog.Logger, eventReceiver EventReceiver, userRegister UserRegister) *Service {
	return &Service{
		log:           log,
		eventReceiver: eventReceiver,
		userRegister:  userRegister,
	}
}

func (s *Service) GetEvents(ctx context.Context) ([]*pb.Event, error) {
	events, err := s.eventReceiver.GetEvents(ctx)
	if err != nil {
		s.log.Error("operation", opGetEvents, "error", err)
		return nil, fmt.Errorf("%s: %w", opGetEvents, err)
	}
	return events, nil
}

func (s *Service) GetEvent(ctx context.Context, eventID string) (*pb.Event, error) {
	event, err := s.eventReceiver.GetEvent(ctx, eventID)
	if err != nil {
		s.log.Error("operation", opGetEvent, "error", err)
		return nil, fmt.Errorf("%s: %w", opGetEvent, err)
	}
	return event, nil
}

func (s *Service) RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error) {
	result, err := s.userRegister.RegisterUser(ctx, eventID, chatID, username)
	if err != nil {
		s.log.Error("operation", opRegisterUser, "error", err)
		return false, fmt.Errorf("%s: %w", opRegisterUser, err)
	}
	return result, nil
}
