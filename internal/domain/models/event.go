package models

import "time"

// Event описывает соответствующую модель данных
type Event struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartsAt    time.Time `db:"starts_at"`
}
