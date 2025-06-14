package main

import (
	"database/sql"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/lulzshadowwalker/personal/internal/config"
	"github.com/lulzshadowwalker/personal/internal/psql"
	"github.com/pressly/goose/v3"
)

func main() {
	db, err := sql.Open("postgres", config.DBConnectionString())
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	goose.SetBaseFS(psql.Migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migration"); err != nil {
		panic(err)
	}
}
