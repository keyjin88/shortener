package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/keyjin88/shortener/internal/app/logger"
	"github.com/keyjin88/shortener/internal/app/storage"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"time"

	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code, err := runMain(m)
	if err != nil {
		logger.Log.Fatal(err)
	}
	os.Exit(code)
}

const (
	testDBName       = "test"
	testUserName     = "test"
	testUserPassword = "test"
)

var (
	getDSN          func() string
	getSUConnection func() (*pgx.Conn, error)
)

func initGetDSN(hostAndPort string) {
	getDSN = func() string {
		return fmt.Sprintf(
			"postgres://%s:%s@%s/%s?sslmode=disable",
			testUserName,
			testUserPassword,
			hostAndPort,
			testDBName,
		)
	}
}

func initGetSUConnection(hostPort string) error {
	host, port, err := getHostPort(hostPort)
	if err != nil {
		return fmt.Errorf("failed to extract the host and port parts from the string %s: %w", hostPort, err)
	}
	getSUConnection = func() (*pgx.Conn, error) {
		conn, err := pgx.Connect(pgx.ConnConfig{
			Host:     host,
			Port:     port,
			Database: "postgres",
			User:     "postgres",
			Password: "postgres",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get a super user connection: %w", err)
		}
		return conn, nil
	}
	return nil
}

func runMain(m *testing.M) (int, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return 1, fmt.Errorf("failed to initialize a pool: %w", err)
	}
	pg, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "15.3",
			Name:       "migrations-integration-tests",
			Env: []string{
				"POSTGRES_USER=postgres",
				"POSTGRES_PASSWORD=postgres",
			},
			ExposedPorts: []string{"5432"},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		return 1, fmt.Errorf("failed to run the postgres container: %w", err)
	}

	defer func() {
		if err := pool.Purge(pg); err != nil {
			logger.Log.Infof("failed to purge the postgres container: %v", err)
		}
	}()

	hostPort := pg.GetHostPort("5432/tcp")
	initGetDSN(hostPort)
	if err := initGetSUConnection(hostPort); err != nil {
		return 1, err
	}

	pool.MaxWait = 10 * time.Second
	var conn *pgx.Conn
	if err := pool.Retry(func() error {
		conn, err = getSUConnection()
		if err != nil {
			return fmt.Errorf("failed to connect to the DB: %w", err)
		}
		return nil
	}); err != nil {
		return 1, fmt.Errorf("failed to retry connect to the DB: %w", err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Log.Infof("failed to correctly close the connection: %v", err)
		}
	}()

	if err := createTestDB(conn); err != nil {
		return 1, fmt.Errorf("failed to crete a test DB: %w", err)
	}

	exitCode := m.Run()
	return exitCode, nil
}

func createTestDB(conn *pgx.Conn) error {
	_, err := conn.Exec(
		fmt.Sprintf(
			`CREATE USER %s PASSWORD '%s'`,
			testUserName,
			testUserPassword,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create a test user: %w", err)
	}

	_, err = conn.Exec(
		fmt.Sprintf(`
			CREATE DATABASE %s
				OWNER '%s'
				ENCODING 'UTF8'
				LC_COLLATE = 'en_US.utf8'
				LC_CTYPE = 'en_US.utf8'
			`, testDBName, testUserName,
		),
	)

	if err != nil {
		return fmt.Errorf("failed to create a test DB: %w", err)
	}

	return nil
}

func getHostPort(hostPort string) (string, uint16, error) {
	hostPortParts := strings.Split(hostPort, ":")
	if len(hostPortParts) != 2 {
		return "", 0, fmt.Errorf("got an invalid host-port string: %s", hostPort)
	}
	portStr := hostPortParts[1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("failed to cast the port %s to an int: %w", portStr, err)
	}
	return hostPortParts[0], uint16(port), nil
}

func initRepository() *URLRepositoryPostgres {
	dsn := getDSN()
	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Log.Errorf("error while initialising DB Pool: %v", err)
		return nil
	}
	ch := make(chan storage.UserURLs)
	repository, err := NewPostgresRepository(dbPool, context.Background(), ch)
	if err != nil {
		logger.Log.Errorf("failed to create new Postgres Repository: %v", err)
	}
	defer repository.Close()
	return repository
}

func TestURLRepositoryPostgres_Save(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	repository := initRepository()

	tests := []struct {
		name               string
		url1               *storage.ShortenedURL
		url2               *storage.ShortenedURL
		ExpectedErrMessage any
	}{
		{
			name: "save success",
			url1: &storage.ShortenedURL{
				UserID:      "userId",
				UUID:        uuid.NewString(),
				ShortURL:    "shortUrl",
				OriginalURL: "Original URL1",
			},
			url2: &storage.ShortenedURL{
				UserID:      "userId",
				UUID:        uuid.NewString(),
				ShortURL:    "shortUrl",
				OriginalURL: "Original URL2",
			},
			ExpectedErrMessage: nil,
		},
		{
			name: "save with error",
			url1: &storage.ShortenedURL{
				UserID:      "userId",
				UUID:        uuid.NewString(),
				ShortURL:    "shortUrl",
				OriginalURL: "Original URL",
			},
			url2: &storage.ShortenedURL{
				UserID:      "userId",
				UUID:        uuid.NewString(),
				ShortURL:    "shortUrl",
				OriginalURL: "Original URL",
			},
			ExpectedErrMessage: "ERROR: duplicate key value violates unique constraint " +
				"\"shortened_url_original_url_key\" (SQLSTATE 23505)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = repository.Save(tt.url1)
			err = repository.Save(tt.url2)
			if err != nil {
				assert.Equal(t, tt.ExpectedErrMessage, err.Error())
			}
		})
	}
}

func TestURLRepositoryPostgres_FindByShortenedURL(t *testing.T) {
	err := logger.Initialize("info")
	if err != nil {
		t.Fatal(err)
	}

	repository := initRepository()

	tests := []struct {
		name                   string
		url1                   *storage.ShortenedURL
		ExpectedSaveErrMessage any
		ExpectedFindErrMessage any
	}{
		{
			name: "find by shortened success",
			url1: &storage.ShortenedURL{
				UserID:      "userId",
				UUID:        uuid.NewString(),
				ShortURL:    "shortUrl",
				OriginalURL: "Original URL1",
			},
			ExpectedSaveErrMessage: nil,
			ExpectedFindErrMessage: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = repository.Save(tt.url1)
			if err != nil {
				assert.Equal(t, tt.ExpectedSaveErrMessage, err.Error())
			}
			url, err2 := repository.FindByShortenedURL(tt.url1.ShortURL)
			assert.Equal(t, tt.ExpectedFindErrMessage, err2)
			assert.Equal(t, tt.url1.ShortURL, url.ShortURL)
		})
	}
}
