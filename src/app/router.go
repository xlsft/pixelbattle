package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xlsft/pixelbattle/middleware"
	authRoutes "github.com/xlsft/pixelbattle/routes/auth"
)

func DefineRouter(app *fiber.App) {

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/", authRoutes.HandlePost)
	auth.Get("/", middleware.AuthMiddleware(), authRoutes.HandleGet)

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})
}
