package psql

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionParams struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func Connect(p ConnectionParams) (*pgxpool.Pool, error) {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(p.Username, p.Password),
		Host:   fmt.Sprintf("%s:%s", p.Host, p.Port),
		Path:   p.Name,
	}

	q := dsn.Query()
	q.Add("sslmode", p.SSLMode)

	dsn.RawQuery = q.Encode()

	pool, err := pgxpool.New(context.Background(), dsn.String())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database because %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database because %w", err)
	}

	slog.Info("connected to database", "host", p.Host, "port", p.Port, "name", p.Name)

	return pool, nil
}
