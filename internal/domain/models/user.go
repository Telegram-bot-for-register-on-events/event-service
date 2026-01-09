package models

// User описывает данные для публикации события регистрации в NATS
type User struct {
	ChatID   int64  `json:"chat_id"`
	Username string `json:"username"`
	EventID  string `json:"event_id"`
}
