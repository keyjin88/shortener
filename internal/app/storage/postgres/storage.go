package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
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
				id           serial
					primary key,
				short_url    varchar,
				original_url varchar,
			
				created_at   date not null,
				updated_at   date not null
			);`
	_, err = dbPool.Exec(ctx, query)
	if err != nil {
		return nil, err
	}
	return &URLRepositoryPostgres{dbPool: dbPool}, nil
}

func (r *URLRepositoryPostgres) FindByShortenedURL(shortURL string) (string, error) {
	ctx := context.Background()
	sql := `SELECT * FROM public.shortened_url WHERE short_url = $1`
	shortenedURL := storage.ShortenedURL{}
	err := r.dbPool.QueryRow(ctx, sql, shortURL).
		Scan(&shortenedURL.ID, &shortenedURL.ShortURL, &shortenedURL.OriginalURL, &shortenedURL.CreatedAt, &shortenedURL.UpdatedAt)
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
		ShortURL:    shortURL, // leave blank for now
		OriginalURL: url,
	}
	sql := `INSERT INTO shortened_url (created_at, updated_at, short_url, original_url) VALUES ($1, $2, $3, $4) RETURNING id;`
	err := r.dbPool.QueryRow(ctx, sql, shortenedURL.CreatedAt, shortenedURL.UpdatedAt, shortenedURL.ShortURL, shortenedURL.OriginalURL).Scan(&shortenedURL.ID)
	if err != nil {
		return storage.ShortenedURL{}, err
	}

	return shortenedURL, nil
}
