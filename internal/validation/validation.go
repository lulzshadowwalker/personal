package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var messages map[string]string = map[string]string{
	"required": "%s is required.",
	"email":    "%s must be a valid email address.",
	"min":      "%s must be at least %d characters long.",
	"max":      "%s must be at most %d characters long.",
	"gte":      "%s must be greater than or equal to %d.",
	"lte":      "%s must be less than or equal to %d.",
	"eq":       "%s must be equal to %d.",
	"ne":       "%s must not be equal to %d.",
	"len":      "%s must be exactly %d characters long.",
	"alpha":    "%s must contain only alphabetic characters.",
	"alphanum": "%s must contain only alphanumeric characters.",
	"numeric":  "%s must contain only numeric characters.",
}

func message(tag string, field string) string {
	if msg, ok := messages[tag]; ok {
		return fmt.Sprintf(msg, field)
	}
	return fmt.Sprintf("%s is invalid.", field)
}

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       any
	}

	XValidator struct {
		validator *validator.Validate
	}

	XValidationError struct {
		Errors []ErrorResponse
	}
)

func (e XValidationError) Error() string {
	var msg string
	for _, err := range e.Errors {
		msg += "[" + err.FailedField + "]: '" + err.Value.(string) + "' | Needs to implement '" + err.Tag + "' and "
	}

	if len(msg) > 0 {
		//  NOTE: Remove the last " and "
		msg = msg[:len(msg)-5]
	}

	return msg
}

func (e XValidationError) Get(field string) string {
	field = strings.ToLower(field)
	var res string
	for _, err := range e.Errors {
		if strings.ToLower(err.FailedField) == field {
			res = message(err.Tag, err.FailedField)
			break
		}
	}

	return res
}

// This is the validator instance
// for more information see: https://github.com/go-playground/validator
var validate = validator.New()

func (v XValidator) Validate(data any) error {
	errors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			errors = append(errors, elem)
		}
	}

	return XValidationError{
		Errors: errors,
	}
}

func NewValidator() *XValidator {
	return &XValidator{
		validator: validate,
	}
}

// func main() {
// 	myValidator := &XValidator{
// 		validator: validate,
// 	}
//
// 	app := fiber.New(fiber.Config{
// 		// Global custom error handler
// 		ErrorHandler: func(c *fiber.Ctx, err error) error {
// 			return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
// 				Success: false,
// 				Message: err.Error(),
// 			})
// 		},
// 	})
//
// 	// Custom struct validation tag format
// 	myValidator.validator.RegisterValidation("teener", func(fl validator.FieldLevel) bool {
// 		// User.Age needs to fit our needs, 12-18 years old.
// 		return fl.Field().Int() >= 12 && fl.Field().Int() <= 18
// 	})
//
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		user := &User{
// 			Name: c.Query("name"),
// 			Age:  c.QueryInt("age"),
// 		}
//
// 		// Validation
// 		if errs := myValidator.Validate(user); len(errs) > 0 && errs[0].Error {
// 			errMsgs := make([]string, 0)
//
// 			for _, err := range errs {
// 				errMsgs = append(errMsgs, fmt.Sprintf(
// 					"[%s]: '%v' | Needs to implement '%s'",
// 					err.FailedField,
// 					err.Value,
// 					err.Tag,
// 				))
// 			}
//
// 			return &fiber.Error{
// 				Code:    fiber.ErrBadRequest.Code,
// 				Message: strings.Join(errMsgs, " and "),
// 			}
// 		}
//
// 		// Logic, validated with success
// 		return c.SendString("Hello, World!")
// 	})
//
// 	log.Fatal(app.Listen(":3000"))
// }

/**
OUTPUT

[1]
Request:

GET http://127.0.0.1:3000/

Response:

{"success":false,"message":"[Name]: '' | Needs to implement 'required' and [Age]: '0' | Needs to implement 'required'"}

[2]
Request:

GET http://127.0.0.1:3000/?name=efdal&age=9

Response:
{"success":false,"message":"[Age]: '9' | Needs to implement 'teener'"}

[3]
Request:

GET http://127.0.0.1:3000/?name=efdal&age=

Response:
{"success":false,"message":"[Age]: '0' | Needs to implement 'required'"}

[4]
Request:

GET http://127.0.0.1:3000/?name=efdal&age=18

Response:
Hello, World!

**/
