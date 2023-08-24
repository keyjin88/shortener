package storage

import "time"

// ShortenURLRequest is a request to shorten the URL
type ShortenURLRequest struct {
	URL string `json:"url"`
}

// ShortenURLResponse is a response to shorten the URL
type ShortenURLResponse struct {
	Result string `json:"result"`
}

// ShortenURLBatchRequest is a request to shorten the URL batch request
type ShortenURLBatchRequest []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ShortenURLBatchResponse is a response to shorten the URL batch
type ShortenURLBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// UsersURLResponse is a response to shorten the URL response for given user ID
type UsersURLResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ShortenedURL is a model for storage
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

// UserURLs is a list of user URLs
type UserURLs struct {
	UserID string
	URLs   []string
}
