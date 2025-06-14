package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/lulzshadowwalker/personal/internal"
	"github.com/lulzshadowwalker/personal/internal/http/session"
	"github.com/lulzshadowwalker/personal/internal/http/template/login"
	"github.com/lulzshadowwalker/personal/internal/validation"
)

type AuthHandler struct {
	s AuthService
}

type AuthService interface {
	Authenticate(ctx context.Context, email string, password string) (internal.User, error)
}

func NewAuth(s AuthService) *AuthHandler {
	return &AuthHandler{s}
}

type LoginRequest struct {
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=8"`
}

func (h *AuthHandler) Index(c *fiber.Ctx) error {
	return Render(c, login.Login())
}

func (h *AuthHandler) Store(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	var errors login.LoginFormErrors
	var params login.LoginFormParams

	if err := validation.NewValidator().Validate(&req); err != nil {
		err := err.(validation.XValidationError)
		//  TODO: There has to better DX to handle this  
		errors.Email = err.Get("email")
		errors.Password = err.Get("password")
	}

	if errors != (login.LoginFormErrors{}) {
		params = login.LoginFormParams{
			Email: req.Email,
		}

		return Render(c, login.LoginForm(params, errors))
	}

	if _, err := h.s.Authenticate(c.Context(), req.Email, req.Password); err != nil {
		s, err := session.New(c)
		if err != nil {
			return err
		}
		defer s.Save()

		s.Errors("email", "Invalid email or password")
		s.Olds("email", req.Email)
		s.Danger("Invalid email or password", "please try again")

		return HxRedirect(c, "/login")
	}

	return HxRedirect(c, "/")
}
