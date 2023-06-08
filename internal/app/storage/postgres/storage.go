package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"time"
)

type URLRepositoryPostgres struct {
	dbPool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool, ctx context.Context) (*URLRepositoryPostgres, error) {
	query := `create table if not exists public.shortened_url
(
    id             serial primary key,
    short_url      varchar,
    original_url   varchar unique,
    created_at     date not null,
    updated_at     date not null,
    correlation_id varchar
);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		return nil, err
	}
	return &URLRepositoryPostgres{dbPool: pool}, nil
}

func (r *URLRepositoryPostgres) FindByShortenedURL(shortURL string) (string, error) {
	ctx := context.Background()
	query := `SELECT original_url FROM public.shortened_url WHERE short_url = $1`
	var originalURL string
	err := r.dbPool.QueryRow(ctx, query, shortURL).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, nil
}

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

func (r *URLRepositoryPostgres) Save(shortenedURL *storage.ShortenedURL) error {
	ctx := context.Background()
	now := time.Now()
	shortenedURL.CreatedAt = now
	shortenedURL.UpdatedAt = now

	query := `INSERT INTO shortened_url (created_at, updated_at, short_url, original_url)
			VALUES ($1, $2, $3, $4)
			RETURNING id;`
	err := r.dbPool.QueryRow(ctx, query, shortenedURL.CreatedAt, shortenedURL.UpdatedAt, shortenedURL.ShortURL, shortenedURL.OriginalURL).Scan(&shortenedURL.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *URLRepositoryPostgres) SaveBatch(urls *[]storage.ShortenedURL) error {
	ctx := context.Background()
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return err
	}

	for _, url := range *urls {
		err = tx.QueryRow(context.Background(),
			`INSERT INTO shortened_url(short_url, original_url, created_at, updated_at, correlation_id) 
				 VALUES ($1, $2, $3, $4, $5) 
				 RETURNING id`,
			url.ShortURL, url.OriginalURL, url.CreatedAt, time.Now(), url.CorrelationID).
			Scan(&url.ID)
		if err != nil {
			return err
		}
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(context.Background())
			logger.Log.Error(p)
		} else if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	return nil
}

func (r *URLRepositoryPostgres) Close() {
	r.dbPool.Close()
}

func (r *URLRepositoryPostgres) Ping() error {
	return r.dbPool.Ping(context.Background())
}
