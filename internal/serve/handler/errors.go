package handler

import (
	"github.com/gofiber/fiber/v2"
)

// CatchAll catches all requests without a corresponding handler.
func CatchAll(c *fiber.Ctx) error {
	return c.
		Status(404).
		SendString("not found")
}
