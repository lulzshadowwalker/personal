package session

import (
	"context"
	"encoding/gob"
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
	"github.com/lulzshadowwalker/personal/internal/config"
	"github.com/lulzshadowwalker/personal/internal/psql"
)

var storage *postgres.Storage
var store *session.Store

type Session struct {
	FiberSession *session.Session
	Flashes      []FlashMessage
	errors       map[string]string
	olds         map[string]string
}

func New(c *fiber.Ctx) (*Session, error) {
	sess, err := store.Get(c)
	if err != nil {
		return nil, err
	}

	s := &Session{
		FiberSession: sess,
		Flashes:      nil,
	}

	flashes, err := getFlashes(c)
	if err != nil {
		return nil, err
	}
	s.Flashes = flashes

	s.errors = make(map[string]string)
	s.olds = make(map[string]string)

	for _, k := range s.FiberSession.Keys() {
		v := s.FiberSession.Get(k)

		if v == nil {
			continue
		}

		switch true {
		case strings.HasPrefix(k, "error."):
			errorKey := k[len("error."):]
			s.errors[errorKey] = fmt.Sprintf("%v", v)
		case strings.HasPrefix(k, "old."):
			oldKey := k[len("old."):]
			s.olds[oldKey] = fmt.Sprintf("%v", v)
		case k == "flashes":
			// Flashes are already handled above
		default:
			slog.Warn("unexpected session key", "key", k, "value", v)
		}
	}

	return s, nil
}

func init() {
	gob.Register(&map[string]string{})
	gob.Register(&[]FlashMessage{})

	pool, err := psql.Connect(psql.ConnectionParams{
		Host:     config.DBHost(),
		Port:     config.DBPort(),
		Username: config.DBUsername(),
		Password: config.DBPassword(),
		Name:     config.DBName(),
		SSLMode:  config.DBSSLMode(),
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	storage = postgres.New(postgres.Config{
		DB:    pool,
		Table: "fiber_storage",
		Reset: false,
	})

	store = session.New(session.Config{Storage: storage})
}

func (s *Session) Save() error {
	for k, v := range s.errors {
		s.FiberSession.Set("error."+k, v)
	}

	for k, v := range s.olds {
		s.FiberSession.Set("old."+k, v)
	}

	s.FiberSession.Set("flashes", s.Flashes)

	if err := s.FiberSession.Save(); err != nil {
		slog.Error("failed to save session", "error", err)
		return err
	}

	return nil
}

func (s *Session) Flush() error {
	return s.FiberSession.Reset()
}

const contextKey = "session"

func FromContext(ctx context.Context) (*Session, error) {
	sess, ok := ctx.Value(contextKey).(*Session)
	if !ok {
		return nil, fmt.Errorf("session not found in context")
	}

	return sess, nil
}

func MustFromContext(ctx context.Context) *Session {
	sess, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}

	return sess
}

func (s *Session) Ingest(ctx context.Context) context.Context {
	return Ingest(ctx, s)
}

func Ingest(ctx context.Context, s *Session) context.Context {
	return context.WithValue(ctx, contextKey, s)
}
