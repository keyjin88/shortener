package postgres

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"time"
)

type URLRepositoryPostgres struct {
	dbPool *pgxpool.Pool
}

func InitPgRepository(ctx context.Context, dataBaseDSN string) (*URLRepositoryPostgres, error) {
	dbPool, err := pgxpool.New(ctx, dataBaseDSN)
	if err != nil {
		return nil, err
	}

	query := `create table if not exists public.shortened_url
(
    id             serial
        primary key,
    short_url      varchar,
    original_url   varchar,

    created_at     date    not null,
    updated_at     date    not null,
    correlation_id varchar
);`
	_, err = dbPool.Exec(ctx, query)
	if err != nil {
		return nil, err
	}
	return &URLRepositoryPostgres{dbPool: dbPool}, nil
}

func (r *URLRepositoryPostgres) FindByShortenedURL(shortURL string) (string, error) {
	ctx := context.Background()
	query := `SELECT * FROM public.shortened_url WHERE short_url = $1`
	shortenedURL := storage.ShortenedURL{}
	var correlationId sql.NullString
	err := r.dbPool.QueryRow(ctx, query, shortURL).
		Scan(&shortenedURL.ID, &shortenedURL.ShortURL, &shortenedURL.OriginalURL, &shortenedURL.CreatedAt, &shortenedURL.UpdatedAt, &correlationId)
	shortenedURL.CorrelationID = correlationId.String
	if err != nil {
		return "", err
	}

	return shortenedURL.OriginalURL, nil
}

func (r *URLRepositoryPostgres) Save(shortURL string, url string) (storage.ShortenedURL, error) {
	ctx := context.Background()
	now := time.Now()
	shortenedURL := storage.ShortenedURL{
		CreatedAt:   now,
		UpdatedAt:   now,
		ShortURL:    shortURL,
		OriginalURL: url,
	}
	query := `INSERT INTO shortened_url (created_at, updated_at, short_url, original_url)
			VALUES ($1, $2, $3, $4)
			RETURNING id;`
	err := r.dbPool.QueryRow(ctx, query, shortenedURL.CreatedAt, shortenedURL.UpdatedAt, shortenedURL.ShortURL, shortenedURL.OriginalURL).Scan(&shortenedURL.ID)
	if err != nil {
		return storage.ShortenedURL{}, err
	}

	return shortenedURL, nil
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
