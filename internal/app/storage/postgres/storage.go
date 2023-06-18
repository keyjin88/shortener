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
    id             serial
        primary key,
    user_id        varchar not null,
    short_url      varchar,
    original_url   varchar
        unique,
    created_at     date    not null,
    updated_at     date    not null,
    correlation_id varchar,
    is_deleted     boolean not null default false
);`
	_, err := pool.Exec(ctx, query)
	if err != nil {
		return nil, err
	}
	return &URLRepositoryPostgres{dbPool: pool}, nil
}

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

func (r *URLRepositoryPostgres) DeleteRecords(ids []string, userId string) error {
	sql := `UPDATE shortened_url SET is_deleted = true WHERE short_url = ANY($1) AND user_id = $2`
	ch := make(chan []string)
	defer close(ch)

	// Запускаем горутину для фан-ин паттерна
	go fanIn(ch, func(ids []string) error {
		// Выполняем множественное обновление
		_, err := r.dbPool.Exec(context.Background(), sql, ids, userId)
		if err != nil {
			logger.Log.Infof("error while deleting: %e", err)
			return err
		}
		return nil
	})
	// Разбиваем слайс на части для более эффективного обновления
	batchSize := 100
	for i := 0; i < len(ids); i += batchSize {
		end := i + batchSize
		if end > len(ids) {
			end = len(ids)
		}
		ch <- ids[i:end]
	}
	return nil
}

func fanIn(ch <-chan []string, f func([]string) error) {
	// Выполняем обновление для каждого блока данных, полученного из канала
	for ids := range ch {
		if err := f(ids); err != nil {
			logger.Log.Infof("error while deleting in fanin: %e", err)
		}
	}
}

func (r *URLRepositoryPostgres) Close() {
	r.dbPool.Close()
}

func (r *URLRepositoryPostgres) Ping(ctx context.Context) error {
	return r.dbPool.Ping(ctx)
}
