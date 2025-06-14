package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/lulzshadowwalker/personal/internal"
	"github.com/lulzshadowwalker/personal/internal/config"
	"github.com/lulzshadowwalker/personal/internal/http/handler"
	"github.com/lulzshadowwalker/personal/internal/http/middleware"
	"github.com/lulzshadowwalker/personal/internal/psql"
	"github.com/lulzshadowwalker/personal/internal/psql/db"
	"github.com/lulzshadowwalker/personal/internal/psql/store"
)

func main() {
	app := fiber.New()
	app.Use(middleware.Session)

	app.Static("/public", "./cmd/http/public")

	app.Get("/", handler.NewHome().Index)
	app.Get("/login", handler.NewAuth(MockAuthService{}).Index)
	app.Post("/login", handler.NewAuth(MockAuthService{}).Store)

	go func() {
		// app.Listen(":" + config.Port())
		if err := app.Listen(":" + config.Port()); err != nil {
			//  TODO: Log and check for server shutdown error
		}
	}()

	slog.Info("server is running", "port", config.Port())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	slog.Info("received shutdown signal, shutting down server gracefully")

	if err := app.Shutdown(); err != nil {
		slog.Error("failed to shutdown server gracefully", "err", err)
		os.Exit(1)
	}

	//  NOTE: Cleanup tasks can be added here, such as closing database connections, etc.
	slog.Info("server shutdown gracefully")
}

type MockAuthService struct {
	//
}

func (m MockAuthService) Authenticate(ctx context.Context, email, password string) (internal.User, error) {
	if email == "email@example.com" && password == "password" {
		return internal.User{}, errors.New("invalid credentials")
	}

	pool, err := psql.Connect(psql.ConnectionParams{
		Host:     config.DBHost(),
		Port:     config.DBPort(),
		Username: config.DBUsername(),
		Password: config.DBPassword(),
		Name:     config.DBName(),
		SSLMode:  config.DBSSLMode(),
	})
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	store := store.NewUserStore(db.New(pool))

	u, _ := store.CreateUser(context.Background(), internal.User{
		Name:  strings.Split(email, "@")[0],
		Email: email,
	})

	return u, nil
}
