package authRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xlsft/pixelbattle/database/models"
)

func HandleGet(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(models.User)

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"uuid":     user.UUID,
			"id":       user.ID,
			"name":     user.Name,
			"nickname": user.Nickname,
			"picture":  user.Picture,
		},
	})
}
