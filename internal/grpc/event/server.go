package grpcserver

import (
	"context"

	"github.com/Telegram-bot-for-register-on-events/shared-proto/pb/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// EventService описывает методы для взаимодействия с сервисным слоем
type EventService interface {
	GetEvents(ctx context.Context) ([]*event.Event, error)
	GetEvent(ctx context.Context, eventID string) (*event.Event, error)
	RegisterUser(ctx context.Context, eventID string, chatID int64, username string) (bool, error)
}

type serverAPI struct {
	event.UnimplementedEventServiceServer
	events EventService
}

// Register регистрирует обработчик, который обрабатывает запросы, приходящие на gRPC-сервер
func Register(grpc *grpc.Server, events EventService) {
	event.RegisterEventServiceServer(grpc, &serverAPI{events: events})
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
	result, err := s.events.RegisterUser(ctx, req.GetEventId(), req.GetChatId(), req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &event.RegisterUserResponse{Success: result}, nil
}
