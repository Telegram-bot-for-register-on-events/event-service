package models

import "time"

// Event описывает соответствующую модель данных
type Event struct {
	ID          string
	Title       string
	Description string
	StartsAt    time.Time
}
