package session

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

//  WARNING:
//  Flash messages currently only work upon redirects and do not work if called in the same handler
//  which I think is the intended behavior following the Post/Redirect/Get pattern.

type FlashMessage struct {
	Type        FlashType
	Title       string
	Description string
}

type FlashType int

const (
	FlashTypeSuccess = iota
	FlashTypeDanger
	FlashTypeInfo
	FlashTypeWarning
	FlashTypePrimary
)

func NewFlashMessage(t FlashType, title, description string) FlashMessage {
	return FlashMessage{
		Type:        t,
		Title:       title,
		Description: description,
	}
}

func (s *Session) Flash(t FlashType, title, description string) {
	flash := NewFlashMessage(t, title, description)
	s.Flashes = append(s.Flashes, flash)
}

func (s *Session) Success(title, description string) {
	s.Flash(FlashTypeSuccess, title, description)
}

func (s *Session) Danger(title, description string) {
	s.Flash(FlashTypeDanger, title, description)
}

func (s *Session) Info(title, description string) {
	s.Flash(FlashTypeInfo, title, description)
}

func (s *Session) Warning(title, description string) {
	s.Flash(FlashTypeWarning, title, description)
}

func (s *Session) Primary(title, description string) {
	s.Flash(FlashTypePrimary, title, description)
}

// GetFlashes retrieves and returns flash messages, then clears them from the session.
func getFlashes(c *fiber.Ctx) ([]FlashMessage, error) {
	sess, err := store.Get(c)
	if err != nil {
		return nil, err
	}

	val := sess.Get("flashes")
	var flashes []FlashMessage
	switch v := val.(type) {
	case nil:
		flashes = make([]FlashMessage, 0)
	case []FlashMessage:
		flashes = v
	case *[]FlashMessage:
		if v != nil {
			flashes = *v
		} else {
			flashes = make([]FlashMessage, 0)
		}
	default:
		panic("unexpected type for flashes: " + fmt.Sprintf("%T", v))
	}

	return flashes, nil
}
