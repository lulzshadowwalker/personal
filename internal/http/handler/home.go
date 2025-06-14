package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lulzshadowwalker/personal/internal/http/template/home"
)

type HomeHandler struct {
	//
}

func NewHome() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Index(c *fiber.Ctx) error {
	//r, _ := session.Primary(w, r, "welcome home", "there")
	// session.Primary(w, r, "welcome home", "there")
	return Render(c, home.Index("lulzie"))
}
