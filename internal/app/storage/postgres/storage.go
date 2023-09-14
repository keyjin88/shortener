package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/pkg/errors"
	"time"
)

// URLRepositoryPostgres is Postgres repository.
type URLRepositoryPostgres struct {
	dbPool       *pgxpool.Pool
	urlsToDelete chan storage.UserURLs
}

// NewPostgresRepository creates a new URLRepositoryPostgres.
func NewPostgresRepository(pool *pgxpool.Pool,
	ctx context.Context,
	toDeleteChan chan storage.UserURLs) (*URLRepositoryPostgres, error) {
	query := `create table if not exists public.shortened_url
(
    id             serial
        primary key,
    user_id        varchar not null,
    short_url      varchar,
    original_url   varchar
        unique,
    created_at     timestamp    not null,
    updated_at     timestamp    not null,
    correlation_id varchar,
    is_deleted     boolean not null default false
);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "error while executing query")
	}
	go WorkerDeleteURLs(toDeleteChan, pool)
	return &URLRepositoryPostgres{dbPool: pool, urlsToDelete: toDeleteChan}, nil
}

// FindByShortenedURL find URL by given shortened string in DB.
func (r *URLRepositoryPostgres) FindByShortenedURL(shortURL string) (storage.ShortenedURL, error) {
	ctx := context.Background()
	query := `SELECT original_url, is_deleted FROM public.shortened_url WHERE short_url = $1`
	var originalURL string
	var isDeleted bool
	err := r.dbPool.QueryRow(ctx, query, shortURL).Scan(&originalURL, &isDeleted)
	if err != nil {
		return storage.ShortenedURL{}, err
	}
	return storage.ShortenedURL{OriginalURL: originalURL, IsDeleted: isDeleted}, nil
}

// FindByOriginalURL find shortened URL by original URL
func (r *URLRepositoryPostgres) FindByOriginalURL(originalURL string) (string, error) {
	ctx := context.Background()
	query := `SELECT short_url FROM public.shortened_url WHERE original_url = $1`
	var shortURL string
	err := r.dbPool.QueryRow(ctx, query, originalURL).Scan(&shortURL)
	if err != nil {
		return "", err
	}
	return shortURL, nil
}

// FindAllByUserID find URLs by user ID
func (r *URLRepositoryPostgres) FindAllByUserID(userID string) ([]storage.UsersURLResponse, error) {
	ctx := context.Background()
	query := `SELECT short_url, original_url FROM shortened_url WHERE user_id = $1`
	rows, err := r.dbPool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userURLs []storage.UsersURLResponse
	for rows.Next() {
		var shortURL, originalURL string
		err := rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, err
		}
		userURLs = append(userURLs, storage.UsersURLResponse{ShortURL: shortURL, OriginalURL: originalURL})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return userURLs, nil
}

// Save method for saving URL in storage
func (r *URLRepositoryPostgres) Save(shortenedURL *storage.ShortenedURL) error {
	ctx := context.Background()
	now := time.Now()
	shortenedURL.CreatedAt = now
	shortenedURL.UpdatedAt = now

	query := `INSERT INTO shortened_url (user_id, created_at, updated_at, short_url, original_url)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id;`
	err := r.dbPool.QueryRow(ctx, query, shortenedURL.UserID, shortenedURL.CreatedAt, shortenedURL.UpdatedAt, shortenedURL.ShortURL, shortenedURL.OriginalURL).Scan(&shortenedURL.ID)
	if err != nil {
		return err
	}

	return nil
}

// SaveBatch saves a batch of USRs to storage
func (r *URLRepositoryPostgres) SaveBatch(urls *[]storage.ShortenedURL) error {
	ctx := context.Background()
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}

	for _, url := range *urls {
		err = tx.QueryRow(context.Background(),
			`INSERT INTO shortened_url(user_id, short_url, original_url, created_at, updated_at, correlation_id) 
				 VALUES ($1, $2, $3, $4, $5, $6) 
				 RETURNING id`,
			url.UserID, url.ShortURL, url.OriginalURL, url.CreatedAt, time.Now(), url.CorrelationID).
			Scan(&url.ID)
		if err != nil {
			return err
		}
	}

	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback(context.Background())
			if err != nil {
				return
			}
			logger.Log.Error(p)
		} else if err != nil {
			err := tx.Rollback(context.Background())
			if err != nil {
				return
			}
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	return nil
}

// Delete method deleted URLs by given IDs
func (r *URLRepositoryPostgres) Delete(ids []string, userID string) error {
	r.urlsToDelete <- storage.UserURLs{UserID: userID, URLs: ids}
	return nil
}

// WorkerDeleteURLs is used to async delete URLs
func WorkerDeleteURLs(ch <-chan storage.UserURLs, pool *pgxpool.Pool) {
	for userUrls := range ch {
		sql := `UPDATE shortened_url SET is_deleted = true, updated_at = $1  WHERE short_url = ANY($2) AND user_id = $3`
		_, err := pool.Exec(context.Background(), sql, time.Now(), userUrls.URLs, userUrls.UserID)
		if err != nil {
			logger.Log.Infof("error while deleting: %e", err)
		}
	}
}

// Close method closes the repository
func (r *URLRepositoryPostgres) Close() {
	r.dbPool.Close()
}

// Ping method pings storage
func (r *URLRepositoryPostgres) Ping(ctx context.Context) error {
	return r.dbPool.Ping(ctx)
}
