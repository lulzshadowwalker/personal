package template

import (
	"context"

	"github.com/lulzshadowwalker/personal/internal/http/session"
)

func Error(ctx context.Context, key string) string {
	s := session.MustFromContext(ctx)
	return s.Error(key)
}

func HasError(ctx context.Context, key string) bool {
	s := session.MustFromContext(ctx)
	return s.HasError(key)
}

func Old(ctx context.Context, key string) string {
	s := session.MustFromContext(ctx)
	return s.Old(key)
}

func HasOld(ctx context.Context, key string) bool {
	s := session.MustFromContext(ctx)
	return s.HasOld(key)
}

func Flashes(ctx context.Context) []session.FlashMessage {
	s := session.MustFromContext(ctx)
	return s.Flashes
}
