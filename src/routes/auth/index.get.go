package authRoutes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xlsft/pixelbattle/database/models"
	"github.com/xlsft/pixelbattle/utils"
)

func HandleGet(ctx *fiber.Ctx) error {
	user, ok := ctx.Locals("user").(models.User)

	if ok == false {
		return ctx.Status(fiber.StatusUnauthorized).JSON(utils.DefineError("Invalid request body, please try again"))
	}

	return ctx.JSON(fiber.Map{
		"data": fiber.Map{
			"user": fiber.Map{
				"uuid": user.UUID,
				"id":   user.ID,
				"name": user.Name,
			},
		},
	})
}
