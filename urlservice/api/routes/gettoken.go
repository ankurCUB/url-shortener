package routes

import (
	"github.com/gofiber/fiber/v2"
)

func GetTokenRange(c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "cannot connect to DB",
	})
}