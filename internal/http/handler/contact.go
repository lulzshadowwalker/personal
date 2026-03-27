package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lulzshadowwalker/personal/internal/http/template/contact"
)

type ContactHandler struct {
	//
}

func NewContact() *ContactHandler {
	return &ContactHandler{}
}

func (h *ContactHandler) Index(c *fiber.Ctx) error {
	//r, _ := session.Primary(w, r, "welcome contact", "there")
	// session.Primary(w, r, "welcome contact", "there")
	return Render(c, contact.Index("lulzie"))
}
