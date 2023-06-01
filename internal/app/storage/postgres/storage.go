package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type URLRepositoryPostgres struct {
	dbPool *pgxpool.Pool
}

func (r *URLRepositoryPostgres) InitPgRepository(ctx context.Context, dataBaseDSN string) (*URLRepositoryPostgres, error) {
	dbPool, err := pgxpool.New(ctx, dataBaseDSN)
	if err != nil {
		if err != nil {
			return nil, err
		}
	}
	query := "CREATE TABLE IF NOT EXISTS shortener.shortened_url (uuid varchar, short_url varchar, original_url varchar);\n"
	_, err = dbPool.Exec(ctx, query)
	if err != nil {
		return nil, err
	}
	return &URLRepositoryPostgres{dbPool: dbPool}, nil
}

func (r *URLRepositoryPostgres) FindByShortenedURL(shortURL string) (string, bool) {
	return "", false
}
func (r *URLRepositoryPostgres) Save(shortURL string, url string) error {
	return nil
}
