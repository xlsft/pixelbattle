package app

import (
	"github.com/gofiber/fiber/v2"
)

func DefineRouter(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
}
