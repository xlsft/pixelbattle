package app

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func StartService() *fiber.App {
	app := fiber.New(fiber.Config{})
	DefineRouter(app)
	port := "8000"
	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}
	app.Listen(":" + port)
	return app
}
