package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Telegram-bot-for-register-on-events/event-service/internal/domain/models"
	"github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EventService описывает методы для взаимодействия с сервисным слоем
type EventService interface {
	GetEvents(ctx context.Context) ([]*event.Event, error)
	GetEvent(ctx context.Context, eventID string) (*event.Event, error)
}

// Registerer описывает метод для передачи данных о регистрации в сервисный слой
type Registerer interface {
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) error
}

// Publisher описывает метод для публикации сообщения в Nats
type Publisher interface {
	PublishMessage(topic string, data []byte) error
}

// serverAPI описывает API для взаимодействия с gRPC-сервером
type serverAPI struct {
	event.UnimplementedEventServiceServer
	events     EventService
	publisher  Publisher
	registerer Registerer
}

// Register регистрирует обработчик, который обрабатывает запросы, приходящие на gRPC-сервер
func Register(grpc *grpc.Server, events EventService, publisher Publisher, registerer Registerer) {
	event.RegisterEventServiceServer(grpc, &serverAPI{events: events, publisher: publisher, registerer: registerer})
}

// GetEvents обрабатывает входящий запрос на получение всех событий
func (s *serverAPI) GetEvents(ctx context.Context, req *event.GetEventsRequest) (*event.GetEventsResponse, error) {
	events, err := s.events.GetEvents(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &event.GetEventsResponse{Events: events}, nil
}

// GetEvent обрабатывает запрос на получение конкретного события
func (s *serverAPI) GetEvent(ctx context.Context, req *event.GetEventRequest) (*event.GetEventResponse, error) {
	e, err := s.events.GetEvent(ctx, req.GetEventId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &event.GetEventResponse{Event: e}, nil
}

// RegisterUser обрабатывает запрос на регистрацию пользователя на конкретное событие
func (s *serverAPI) RegisterUser(ctx context.Context, req *event.RegisterUserRequest) (*event.RegisterUserResponse, error) {
	err := s.registerer.RegisterUser(ctx, req.GetEventId(), req.GetChatId(), req.GetUsername())
	if err != nil {
		return &event.RegisterUserResponse{Success: false}, fmt.Errorf("events.RegisterUser: %w", err)
	}

	// Формируем сообщение для публикации в шину данных
	user := &models.User{
		ChatID:   req.GetChatId(),
		Username: req.GetUsername(),
		EventID:  req.GetEventId(),
	}

	// Сериализируем данные
	jsonData, err := json.Marshal(user)
	if err != nil {
		return &event.RegisterUserResponse{Success: false}, fmt.Errorf("events.RegisterUser: %w", err)
	}

	// Публикуем сообщение
	err = s.publisher.PublishMessage("register.user", jsonData)
	if err != nil {
		return &event.RegisterUserResponse{Success: false}, fmt.Errorf("events.RegisterUser: %w", err)
	}
	return &event.RegisterUserResponse{Success: true}, nil
}
