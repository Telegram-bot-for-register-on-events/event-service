package models

import "time"

// Registration описывает данные для регистрации пользователя на событие
type Registration struct {
	ID        string
	EventID   string
	ChatID    int64
	Username  string
	CreatedAt time.Time
}
