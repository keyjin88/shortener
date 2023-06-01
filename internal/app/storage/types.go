package storage

import "time"

type ShortenedURL struct {
	ID          int64     `json:"-" db:"id"`
	UUID        string    `json:"uuid" db:"-"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
	ShortURL    string    `json:"short_url" db:"short_url"`
	OriginalURL string    `json:"original_url" db:"original_url"`
}
