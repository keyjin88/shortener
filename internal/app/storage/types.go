package storage

import "time"

type ShortenURLRequest struct {
	URL string `json:"url"`
}

type ShortenURLResponse struct {
	Result string `json:"result"`
}

type ShortenURLBatchRequest []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ShortenURLBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type UsersURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ShortenedURL struct {
	ID            int64     `json:"-" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	UUID          string    `json:"uuid" db:"-"`
	CreatedAt     time.Time `json:"-" db:"created_at"`
	UpdatedAt     time.Time `json:"-" db:"updated_at"`
	ShortURL      string    `json:"short_url" db:"short_url"`
	OriginalURL   string    `json:"original_url" db:"original_url"`
	CorrelationID string    `json:"correlation_id" db:"correlation_id"`
	IsDeleted     bool      `json:"is_deleted" db:"is_deleted"`
}

type UserURLs struct {
	UserID string
	URLs   []string
}
