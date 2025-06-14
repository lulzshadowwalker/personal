package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/lulzshadowwalker/personal/internal/http/session"
)

func HxRedirect(c *fiber.Ctx, to string) error {
	if _, ok := c.GetReqHeaders()["Hx-Request"]; ok {
		c.Set("Hx-Redirect", to)
		c.Set("Hx-Trigger", "redirect")
		return nil
	}

	return c.Redirect(to, http.StatusSeeOther)
}

func Render(c *fiber.Ctx, component templ.Component) error {
	s, err := session.New(c)
	if err != nil {
		return err
	}

	c.SetUserContext(session.Ingest(c.UserContext(), s))

	//  NOTE: Flush the session to ensure that session data
	//  is cleared after being passed to the template
	s.Flush()

	c.Set("Content-Type", "text/html")
	return component.Render(c.UserContext(), c.Response().BodyWriter())
}
