package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lulzshadowwalker/personal/internal/http/session"
)

func Session(c *fiber.Ctx) error {
	s, err := session.New(c)
	if err != nil {
		return err
	}

	c.SetUserContext(session.Ingest(c.UserContext(), s))
	return c.Next()
}
