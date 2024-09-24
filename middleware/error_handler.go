package middleware

import (
	"github.com/gofiber/fiber/v2"
)

type ErrorHandlerMiddleware struct {
	app *fiber.App
}

func NewErrorHandlerMiddleware(app *fiber.App) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		app: app,
	}
}

func (eh *ErrorHandlerMiddleware) ErrorHandler(c *fiber.Ctx) error {
	err := c.Next()

	if err != nil {
		return err
	}

	return nil
}

func (eh *ErrorHandlerMiddleware) Bind() {
	eh.app.Use("/api", eh.ErrorHandler)
}
