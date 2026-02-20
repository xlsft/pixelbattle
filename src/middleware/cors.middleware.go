package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func CorsMiddleware(ctx *fiber.Ctx) error {
	ctx.Set("Access-Control-Allow-Origin", "*")
	ctx.Set("Access-Control-Allow-Methods", "*")
	ctx.Set("Access-Control-Allow-Headers", "*")

	if ctx.Method() == "OPTIONS" {
		return ctx.SendStatus(fiber.StatusNoContent)
	}

	return ctx.Next()
}
